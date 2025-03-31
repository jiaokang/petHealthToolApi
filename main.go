package main

import (
	"fmt"
	"petHealthToolApi/config"
	"petHealthToolApi/core"
	"petHealthToolApi/global"
	"petHealthToolApi/router"
)

// 启动程序
func main() {
	global.Log = core.NewLog()

	// 初始化路由
	router := router.InitRouter()
	address := fmt.Sprintf("%s:%d", config.Config.System.Host, config.Config.System.Port)
	global.Log.Infof("系统启动,运行在:%s", address)
	err := router.Run(address)
	if err != nil {
		panic(err)
	}
}
