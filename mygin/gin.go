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

	"dsjk.com/dwggo/core"
	"dsjk.com/dwggo/lib/helper"
	"dsjk.com/dwggo/lib/helper/str"
	"dsjk.com/dwggo/lib/helper/str/verify"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var pid_file string       //save PID for stop
var limiter *rate.Limiter //限流器

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

func Start() {
	defer core.ForPanicLog()
	defer closeAll()

	port := core.GetConfigString("app.port")
	r := ginInit()
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

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

func ginInit() *gin.Engine {
	gin.SetMode(core.GetConfigString("env"))
	gin.DefaultWriter = core.GetLogIoWriter(core.GetConfigString("logger.gin_path"))
	r := gin.Default()
	//r.SetTrustedProxies([]string{"负载均衡,代理IP"})
	middleWareLimiter(r) //限流器
	InitRoute(r)         //映射路由到控制器

	//r.Delims("${", "}") //默认的{{}}和vue冲突,代码必须在LoadHTMLGlob前面
	r.LoadHTMLGlob("view/**/*")
	return r
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
