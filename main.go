package main

import (
	"fmt"
	"goWeb/conf"
	"goWeb/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 載入設定檔
	settings := conf.GetConfig()
	fmt.Printf("%v", settings.Username)
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 加載靜態頁面
	ginServer.LoadHTMLGlob("views/*")
	// 加載資源文件
	ginServer.Static("/static", "./static")
	// 載入 Router
	routes.SetRouter(ginServer)
	// 啟動
	ginServer.Run(settings.ServerPort)
}
