package main

import (
	"log"
	"os"
	"strings"

	"dsjk.com/openplatform/mygin"
)

func main() {
	cmd() //start/stop/restart
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

	switch strings.ToLower(os_args[0]) {
	case "start":
		log.Println("start ing...")
		mygin.Start()
	case "stop":
		log.Println("stop ing...")
		mygin.Stop()
	case "restart":
		log.Println("restart ing...")
		mygin.Stop()
		mygin.Start()

	default:
		log.Println("cmd arg is wrong,such as: xxx start/stop/restart")
	}
}
