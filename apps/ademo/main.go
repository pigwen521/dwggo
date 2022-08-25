package main

import (
	"log"
	"os"
	"strings"

	"dsjk.com/dwggo/apps/ademo/controller"
	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/mygin"
	"github.com/gin-gonic/gin"
)

func main() {
	cmd() //go run main.go start/stop/restart
}

/**
运行命令
本地开发：go run main.go start/stop/restart
编译后：xxx start/stop/restart
*/
func cmd() {
	os_args := os.Args
	if len(os_args) == 0 {
		log.Println("cmd miss arg,such as: xxx start/stop/restart")
		return
	}
	if len(os_args) == 2 {
		os_args[0] = os_args[1] //暂时只支持一个参数，兼容：go run main.go start 两个参数
	}

	gin := initGin()
	switch strings.ToLower(os_args[0]) {
	case "start":
		log.Println("start ing...")
		mygin.Start(gin, controller.InitCtrlByName)
	case "stop":
		log.Println("stop ing...")
		mygin.Stop()
	case "restart":
		log.Println("restart ing...")
		mygin.ReStart(gin, controller.InitCtrlByName)

	default:
		log.Println("cmd arg is wrong,such as: xxx start/stop/restart")
	}
}

func initGin() *gin.Engine {
	gin.SetMode(core.GetConfigString("env"))
	gin.DefaultWriter = core.GetLogIoWriter(core.GetConfigString("logger.gin_path"))
	r := gin.Default()
	//r.SetTrustedProxies([]string{"负载均衡,代理IP"})

	//r.Delims("${", "}") //默认的{{}}和vue冲突,代码必须在LoadHTMLGlob前面
	r.LoadHTMLGlob("view/**/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	return r
}
