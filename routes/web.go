package routes

import (
	"goWeb/controllers"
	"goWeb/middleware"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func GetWebRouters(r *gin.Engine) {
	// 載入 we共用模板設定
	r.HTMLRender = createMyRender()
	r.GET("/login", controllers.Login)   // 登入頁
	r.POST("/auth", controllers.Auth)    // 登入驗證
	r.GET("/logout", controllers.Logout) // 登出頁
	// 首頁
	r.GET("/", middleware.Auth(), controllers.Index)
	r.GET("/index", middleware.Auth(), controllers.Index)
	// 產品管理
	r.GET("/product", middleware.Auth(), controllers.ProductManage)
}

/*
* 共用模板設定
* @return multitemplate.Renderer
 */
func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	common := map[string]string{
		"menu":   "views/layout/menu.html",
		"header": "views/layout/header.html",
		"layout": "views/layout/layout.html",
	}
	includes := map[string]string{
		"productModal": "views/includes/productModal.html",
	}
	// 首頁
	r.AddFromFiles("index", common["layout"], common["header"], common["menu"], "views/main/index.html")
	// 登入頁
	r.AddFromFiles("login", common["layout"], "views/main/login.html")
	// 存貨管理
	r.AddFromFiles("product", common["layout"], common["header"], common["menu"], includes["productModal"], "views/main/product.html")
	// 404 page
	r.AddFromFiles("404", "views/main/404.html")
	return r
}
