package mygin

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper/arrmap"
	"dsjk.com/dwggo/system/lib/helper/str"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter //限流器
var MyginEngin *gin.Engine

func Stop() *process.Process {
	old_pid, old_proc, err := getOldPidProc()
	if err != nil {
		cmdLog("stop error: pid found error " + err.Error())
		return nil
	}
	if old_pid == 0 {
		cmdLog("stop skip: old pid is not found")
		return nil
	}
	err = old_proc.Terminate()
	if err != nil {
		cmdLog("stop error: kill error,pid:" + str.ToString(old_pid) + ",err:" + err.Error())
		return nil
	} else {
		//cmdLog("stop signal send success")
		return old_proc
	}
}

func ReStart(gin *gin.Engine, initCtrlByNameCB InitCtrlByNameCB) {
	old_proc := Stop()
	err := waitStop(old_proc, 3)
	if err != nil {
		cmdLog("waitStop error:" + err.Error())
		cmdLog("try start...")
	}
	Start(gin, initCtrlByNameCB)
}

// 等待先前启动的进程停止了再start
func waitStop(old_proc *process.Process, waitSecond int) error {
	if old_proc == nil {
		return nil
	}
	_, old_proc, err := getProcessByPid(old_proc.Pid)
	if err != nil {
		return err
	}
	if old_proc == nil { //进程消失了
		return nil
	}

	old_status, err := old_proc.Status()
	if err != nil {
		return errors.New("old process status error " + err.Error())
	}
	if !arrmap.InArrayStr(old_status, []string{"R", "S", "D"}) { //R运行，S中断，D不可中断，Z僵死，T停止
		return nil //算等待好了
	}
	if waitSecond == 0 { //等待超时了
		return errors.New("stop timeout")
	}
	//继续等待
	//if arrmap.InArrayStr(old_status, []string{"R", "S", "D"}) {
	cmdLog("waitStop ... process status: " + old_status)
	time.Sleep(time.Second)
	return waitStop(old_proc, waitSecond-1)
	//}
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
	pid := getNowPid()

	var err error
	if core.GetConfigString("app.scheme") == "https" { //HTTPS证书
		certFile := core.GetConfigString("ssl.certFile")
		keyFile := core.GetConfigString("ssl.keyFile")
		if certFile == "" || keyFile == "keyFile" {
			core.LogErrorAndPanic("https时certFile,keyFile不能为空")
		}
		cmdLog("Listening and serving HTTPS on " + srv.Addr + ",pid:" + str.ToString(pid))
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		cmdLog("Listening and serving HTTP on " + srv.Addr + ",pid:" + str.ToString(pid))
		err = srv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		cmdLog("start error:" + err.Error())
		os.Exit(0)
	}
	cmdLog("old server listenAndServe end,pid:" + str.ToString(pid))
}
func cmdLog(str string) {
	core.LogInfoCustomFmt(str, "cmd")
}

var now_pid int32 //当前启动程序的PID
//var old_pid int32 //先前启动程序的PID

func getNowPid() int32 {
	if now_pid != 0 {
		return now_pid
	}
	now_pid = (int32)(os.Getpid())
	return now_pid
}

func getOldPidProc() (int32, *process.Process, error) {
	//先获取当前进程
	now_proc_name, _, err := getProcessByPid(getNowPid())
	if err != nil {
		return 0, nil, err
	}

	//再根据端口号获取先前进程
	old_pid, err := getPidByPort(core.GetConfigInt("app.port"))
	if err != nil {
		return 0, nil, err
	}
	if old_pid == 0 { //先前还没启动
		return 0, nil, nil
	}

	old_proc_name, old_proc, err := getProcessByPid(old_pid)
	if err != nil {
		return 0, nil, err
	}
	if old_proc == nil { //不存在
		return 0, nil, nil
	}

	//当前和先前进程的进程名一致，才确保PID正确，才能被kill
	if old_proc_name == now_proc_name {
		return old_pid, old_proc, nil
	} else {
		cmdLog("get the old_pid error,old_proc_name!=now_proc_name," + old_proc_name + "!=" + now_proc_name)
		return 0, nil, nil
	}
}

func getPidByPort(port_in int) (int32, error) {
	port := (uint32)(port_in)
	net_conns, err := net.Connections("tcp")
	//{"fd":130,"family":2,"type":1,"localaddr":{"ip":"10.0.3.163","port":63503},"remoteaddr":{"ip":"110.43.121.209","port":443},"status":"CLOSE_WAIT","uids":null,"pid":65919}
	if err != nil {
		return 0, err
	}
	var find_net net.ConnectionStat
	for _, cs := range net_conns {
		if cs.Laddr.Port == port {
			find_net = cs
		}
	}
	if find_net.Pid == 0 {
		return 0, nil
	}
	return find_net.Pid, nil
}
func getProcessByName(name string) *process.Process {
	pids, _ := process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)
		pn_name, _ := pn.Name()
		if pn_name == name {
			return pn
		}
	}
	return nil
}
func getProcessByPid(pid int32) (string, *process.Process, error) {
	now_pn, err := process.NewProcess(pid)
	if errors.Is(err, process.ErrorProcessNotRunning) {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	now_pName, err := now_pn.Name()
	return now_pName, now_pn, err
}

// graceful  优雅关闭
func forShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //退出终端=》syscall.SIGHUP；ctrl+c=>os.Interrt；kill xx =>syscall.SIGTERM,Terminate;；kill -9=>os.Kill syscall.SIGKILL 忽略 无法被接收
	got := <-quit

	cmdLog("old server shutdown start:" + got.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		cmdLog("old Server shutdown error :" + err.Error())
	} else {
		cmdLog("old server shutdown success ^_^,oldpid:" + str.ToString(getNowPid()))
	}
	/* empty_pid := ""
	helper.WriteFile(pid_file, &empty_pid) */
}

func ginInit(r *gin.Engine, initCtrlByNameCB InitCtrlByNameCB) *gin.Engine {
	//middleWareTimeout(r)	//请求执行超时告警
	middleWareLimiter(r)           //限流器
	InitRoute(r, initCtrlByNameCB) //映射路由到控制器
	MyginEngin = r
	return r
}

// 判断模板文件是否存在
// if mygin.IsExistTemplateFile("home/index.html") )
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

// 限流器中间件
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

// 请求超时告警
func middleWareTimeout(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		go func(c *gin.Context) {
			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded { //超时了
					core.LogError("请注意，请求执行超时了" + c.Request.RequestURI)
				}
			}
		}(c)
		c.Next() //执行业务代码
	})
}
func closeAll() {
	core.RedisPoolClose()
}
