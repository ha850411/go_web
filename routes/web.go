package routes

import (
	"goWeb/controllers"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func GetWebRouters(r *gin.Engine) {
	// 載入 we共用模板設定
	r.HTMLRender = createMyRender()
	r.GET("/", controllers.Index)
	r.GET("/index", controllers.Index)
	r.GET("/login", controllers.Login)
	r.GET("/product", controllers.ProductManage)
}

/*
* 共用模板設定
* @return multitemplate.Renderer
 */
func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	menu, header, layout := "views/layout/menu.html", "views/layout/header.html", "views/layout/layout.html"
	// 首頁
	r.AddFromFiles("index", layout, header, menu, "views/main/index.html")
	// 登入頁
	r.AddFromFiles("login", layout, "views/main/login.html")
	// 存貨管理
	r.AddFromFiles("product", layout, header, menu, "views/main/product.html")
	// 404 page
	r.AddFromFiles("404", "views/main/404.html")
	return r
}
