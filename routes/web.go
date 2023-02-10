package routes

import (
	"goWeb/controllers"
	"goWeb/middleware"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func GetWebRouters(r *gin.Engine) {
	// 載入共用模板設定
	r.HTMLRender = createMyRender()
	r.GET("/login", controllers.Login)   // 登入頁
	r.POST("/auth", controllers.Auth)    // 登入驗證
	r.GET("/logout", controllers.Logout) // 登出頁
	// 存貨管理
	r.GET("/", middleware.Auth(), controllers.ProductManage)
	r.GET("/product", middleware.Auth(), controllers.ProductManage)
	// 存貨分析
	r.GET("/analysis", middleware.Auth(), controllers.Analysis)
	// 存貨分析
	r.GET("/order", middleware.Auth(), controllers.OrderManage)
	// 存貨分析
	r.GET("/setting", middleware.Auth(), controllers.SettingManage)

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
		"orderModal":   "views/includes/orderModal.html",
	}
	// 存貨分析
	r.AddFromFiles("analysis", common["layout"], common["header"], common["menu"], "views/main/analysis.html")
	// 登入頁
	r.AddFromFiles("login", common["layout"], "views/main/login.html")
	// 存貨管理
	r.AddFromFiles("product", common["layout"], common["header"], common["menu"], includes["productModal"], "views/main/product.html")
	// 訂單管理
	r.AddFromFiles("order", common["layout"], common["header"], common["menu"], includes["orderModal"], "views/main/order.html")
	// 設定
	r.AddFromFiles("setting", common["layout"], common["header"], common["menu"], includes["productModal"], "views/main/setting.html")
	// 404 page
	r.AddFromFiles("404", "views/main/404.html")
	return r
}
