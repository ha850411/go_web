package main

import (
	"fmt"
	"goWeb/conf"
	"goWeb/routes"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

/*
* 載入 log 設定
 */
func init() {
	var logDir = "./logs"
	_, err := os.Stat(logDir)
	if err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}
	file, _ := os.Create(fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02")))
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}

func main() {
	// 載入設定檔
	settings := conf.GetConfig()
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 載入 Router
	routes.SetRouter(ginServer)
	// 啟動
	ginServer.Run(settings.Server.Port)
}
