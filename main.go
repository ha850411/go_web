package main

import (
	"goWeb/conf"
	"goWeb/routes"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now().Format("2006-01-02 15:04:05")
		c.Set("Test", t)
		c.Next()
	}
}

func main() {

	// 載入設定檔
	settings := conf.GetConfig()
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 中間件
	ginServer.Use(Logger())
	// 載入 Router
	routes.SetRouter(ginServer)

	// 啟動
	ginServer.Run(settings.ServerPort)
}
