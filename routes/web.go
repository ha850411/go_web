package routes

import (
	"goWeb/controllers"
	"goWeb/controllers/front"
	"goWeb/middleware"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func GetWebRouters(r *gin.Engine) {
	// 載入共用模板設定
	r.HTMLRender = createMyRender()
	// ========= 後台 ===========
	adminGroup := r.Group("/admin")
	adminGroup.GET("/login", controllers.Login)   // 登入頁
	adminGroup.POST("/auth", controllers.Auth)    // 登入驗證
	adminGroup.GET("/logout", controllers.Logout) // 登出頁
	// 存貨管理
	adminGroup.GET("/", middleware.Auth(), controllers.ProductManage)
	adminGroup.GET("/product", middleware.Auth(), controllers.ProductManage)
	// 圖片管理
	adminGroup.GET("/product/pic/:pid", middleware.Auth(), controllers.ProductPicManage)
	// 產品類別管理
	adminGroup.GET("/product/type", middleware.Auth(), controllers.ProductTypeManage)
	// 存貨分析
	adminGroup.GET("/analysis", middleware.Auth(), controllers.Analysis)
	// 訂單管理
	adminGroup.GET("/order", middleware.Auth(), controllers.OrderManage)
	// 設定
	adminGroup.GET("/setting", middleware.Auth(), controllers.SettingManage)
	// 變更密碼
	adminGroup.POST("/setting/updatePassword", middleware.Auth(), controllers.UpdatePassword)
	// === line bot ===
	linebotGroup := r.Group("/linebot")
	linebotGroup.GET("/notify", middleware.Auth(), controllers.LineBotNotify)
	// ========= 前台 ===========
	r.GET("/", front.Index)
	r.GET("/product", front.Product)
	r.GET("/product/:id", front.ProductDetail)
	r.GET("/cart", front.ShoppingCart)
	r.GET("/about", front.About)
	r.GET("/contact", front.Contact)
	r.GET("/orders", front.Orders)
	r.POST("/orders/add", front.OrdersAdd)
}

/*
* 共用模板設定
* @return multitemplate.Renderer
 */
func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	common := map[string]string{
		"menu":         "views/layout/menu.html",
		"header":       "views/layout/header.html",
		"layout":       "views/layout/layout.html",
		"front-layout": "views/front/front-layout.html",
		"front-header": "views/front/front-header.html",
	}
	includes := map[string]string{
		"productModal":      "views/includes/productModal.html",
		"orderModal":        "views/includes/orderModal.html",
		"productPicModal":   "views/includes/productPicModal.html",
		"productTypeModal":  "views/includes/productTypeModal.html",
		"front-order-modal": "views/includes/frontOrderModal.html",
	}
	index := map[string]string{
		"layout": "views/index/layout.html",
		"header": "views/index/header.html",
		"banner": "views/index/banner.html",
		"footer": "views/index/footer.html",
	}
	// 存貨分析
	r.AddFromFiles("analysis", common["layout"], common["header"], common["menu"], "views/main/analysis.html")
	// 登入頁
	r.AddFromFiles("login", common["layout"], "views/main/login.html")
	// 存貨管理
	r.AddFromFiles("product", common["layout"], common["header"], common["menu"], includes["productModal"], "views/main/product.html")
	// 圖片管理
	r.AddFromFiles("product.picture", common["layout"], common["header"], common["menu"], includes["productPicModal"], "views/main/productPic.html")
	// 產品類別管理
	r.AddFromFiles("product.type", common["layout"], common["header"], common["menu"], includes["productTypeModal"], "views/main/productType.html")
	// 訂單管理
	r.AddFromFiles("order", common["layout"], common["header"], common["menu"], includes["orderModal"], "views/main/order.html")
	// 設定
	r.AddFromFiles("setting", common["layout"], common["header"], common["menu"], includes["productModal"], "views/main/setting.html")
	// === 前端 ===
	r.AddFromFiles("front-order", common["front-layout"], common["front-header"], includes["front-order-modal"], "views/front/order.html")
	r.AddFromFiles("front-test", common["front-layout"], common["front-header"], "views/front/test.html")
	// 首頁
	r.AddFromFiles("index", index["layout"], index["header"], index["banner"], index["footer"], "views/index/index.html")
	// 所有商品頁
	r.AddFromFiles("front_product", index["layout"], index["header"], index["footer"], "views/index/product.html")
	// 商品資訊
	r.AddFromFiles("front_product_detail", index["layout"], index["header"], index["footer"], "views/index/product-detail.html")
	r.AddFromFiles("front_cart", index["layout"], index["header"], index["footer"], "views/index/cart.html")
	r.AddFromFiles("front_about", index["layout"], index["header"], index["footer"], "views/index/about.html")
	r.AddFromFiles("front_contact", index["layout"], index["header"], index["footer"], "views/index/contact.html")

	// === 404 page ===
	r.AddFromFiles("404", "views/main/404.html")
	return r
}
