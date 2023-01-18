package main

import (
	"goWeb/conf"
	"goWeb/routes"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	menu, header, layout := "views/layout/menu.html", "views/layout/header.html", "views/layout/layout.html"
	// 首頁
	r.AddFromFiles("index", layout, header, menu, "views/main/index.html")
	// 存貨管理
	r.AddFromFiles("product", layout, header, menu, "views/main/product.html")
	// 404
	r.AddFromFiles("404", "views/main/404.html")
	return r
}

func main() {
	// 載入設定檔
	settings := conf.GetConfig()
	// 初始化 gin 物件
	ginServer := gin.Default()
	// 加載靜態頁面
	// ginServer.LoadHTMLGlob("views/**/*")
	ginServer.HTMLRender = createMyRender()
	// 加載資源文件
	ginServer.Static("/static", "./static")
	// 載入 Router
	routes.SetRouter(ginServer)
	// 啟動
	ginServer.Run(settings.ServerPort)
}
