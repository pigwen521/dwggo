package mygin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper"
	"dsjk.com/dwggo/system/lib/helper/str"
	"dsjk.com/dwggo/system/lib/helper/str/verify"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/time/rate"
)

var pid_file string       //save PID for stop
var limiter *rate.Limiter //限流器
var MyginEngin *gin.Engine

func init() {
	port := core.GetConfigString("app.port")
	pid_file = core.GetConfigString("logger.custom_dir") + "gin_" + port + ".pid"
}

func Stop() {
	old_pid, err := getOldPid()
	if err != nil {
		cmdLog("stop error: pid read error " + err.Error() + ",pid file:" + pid_file)
		return
	}
	if old_pid == "" {
		cmdLog("stop error: pid is not found")
		return
	}
	if !verify.IsNumber(&old_pid) { //防止pid文件被注入
		cmdLog("stop error: pid is not number")
		return
	}
	out, err := exec.Command("sh", "-c", "kill "+old_pid).Output()
	if err != nil {
		cmdLog("stop error: cmd exec error,pid:" + old_pid + ",err:" + err.Error())
	} else {
		cmdLog("stop success:" + string(out))
		empty_pid := ""
		helper.WriteFile(pid_file, &empty_pid)
	}
}

func Start(gin *gin.Engine, initCtrlByNameCB InitCtrlByNameCB) {
	defer core.ForPanicLog()
	defer closeAll()

	port := core.GetConfigString("app.port")
	r := ginInit(gin, initCtrlByNameCB)
	//r.Run(":" + port)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go listenAndServer(srv, port) //启动服务
	forShutdown(srv)              //优雅关闭
}
func listenAndServer(srv *http.Server, port string) {
	pid := getPid()
	cmdLog("Listening and serving HTTP on " + srv.Addr + ",pid:" + pid)
	old_pid, _ := getOldPid()
	helper.WriteFile(pid_file, &pid)

	var err error
	if core.GetConfigString("app.scheme") == "https" { //HTTPS证书
		certFile := core.GetConfigString("ssl.certFile")
		keyFile := core.GetConfigString("ssl.keyFile")
		if certFile == "" || keyFile == "keyFile" {
			core.LogErrorAndPanic("https时certFile,keyFile不能为空")
		}
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}
	
	if err != nil && err != http.ErrServerClosed {
		cmdLog("start error:" + err.Error())
		helper.WriteFile(pid_file, &old_pid) //启动失败，写回去
		os.Exit(0)
	}
	cmdLog("ListenAndServe end,pid:" + pid)
}
func cmdLog(str string) {
	core.LogInfoCustomFmt(str, "cmd")
}
func getPid() string {
	return strconv.Itoa(os.Getpid())
}
func getOldPid() (string, error) {
	pid, err := ioutil.ReadFile(pid_file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(pid), err
}

//graceful  优雅关闭
func forShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM) //ctrl+c=>os.Interrt  kill -9=>os.Kill ,kill xx =>syscall.SIGTERM
	got := <-quit

	cmdLog("server shutdown start:" + got.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		cmdLog("Server Shutdown error:" + err.Error())
	} else {
		cmdLog("server shutdown success,oldpid:" + getPid())
	}
}

func ginInit(r *gin.Engine, initCtrlByNameCB InitCtrlByNameCB) *gin.Engine {

	middleWareLimiter(r)           //限流器
	InitRoute(r, initCtrlByNameCB) //映射路由到控制器
	MyginEngin = r
	return r
}

//判断模板文件是否存在
//if mygin.IsExistTemplateFile("home/index.html") )
func IsExistTemplateFile(tmp_file string) bool {
	if gin.IsDebugging() {
		aa := (MyginEngin.HTMLRender).(render.HTMLDebug)
		ren := aa.Instance(tmp_file, nil).(render.HTML)
		return ren.Template.Lookup(tmp_file) != nil
	} else {
		aa := (MyginEngin.HTMLRender).(render.HTMLProduction)
		return aa.Template.Lookup(tmp_file) != nil
	}
}

//限流器中间件
func middleWareLimiter(r *gin.Engine) {
	rate_second := core.GetConfigFloat64("limiter.rate_second")
	capacity := core.GetConfigInt("limiter.capacity")
	if rate_second == 0 || capacity == 0 {
		return
	}
	limiter = rate.NewLimiter(rate.Limit(rate_second), capacity)

	r.Use(func(ctx *gin.Context) {
		if limiter.Allow() == false {
			http.Error(ctx.Writer, http.StatusText(429), http.StatusTooManyRequests)
			core.LogError("当前请求过多，当前限流：" + str.ToString(rate_second) + "次/秒,容量：" + str.ToString(capacity))
			ctx.Abort()
		}
	})
}
func closeAll() {
	core.RedisPoolClose()
}
