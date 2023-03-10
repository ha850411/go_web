package main

import (
	"goWeb/conf"
	"goWeb/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 載入 Router
	routes.SetRouter(ginServer)
	// 啟動
	go ginServer.RunTLS(":443", "./certs/server.crt", "./certs/server.key")
	ginServer.Run(conf.Settings.Server.Port)
}
