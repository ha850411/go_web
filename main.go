package main

import (
	"goWeb/conf"
	"goWeb/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// release mode
	gin.SetMode(gin.ReleaseMode)
	// 載入設定檔
	settings := conf.GetConfig()
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 載入 Router
	routes.SetRouter(ginServer)
	// 啟動
	ginServer.Run(settings.Server.Port)
}
