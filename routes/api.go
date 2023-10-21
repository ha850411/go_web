package routes

import (
	"goWeb/controllers/api"
	"goWeb/controllers/api/front"
	"goWeb/middleware"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	// === products Groups ===
	productsGroup := apiGroup.Group("/products")
	// 取得單一商品資訊
	productsGroup.GET("/:id", api.GetProductInfo)
	// 取得商品列表
	productsGroup.GET("/getList", api.GetProductsList)
	// 更新數量
	productsGroup.PUT("/updateAmount", api.UpdateAmount)
	// 新增產品
	productsGroup.POST("/add", api.AddProduct)
	// 編輯產品
	productsGroup.PATCH("/edit", api.EditProduct)
	// 刪除產品
	productsGroup.DELETE("/delete/:id", api.DeleteProduct)
	// 取得存貨警戒商品
	productsGroup.GET("/getTips", api.GetTips)
	// 取得商品名稱
	productsGroup.GET("/getProductsNameList", api.GetProductsNameList)
	// 取得商品 log
	productsGroup.GET("/getProductsLog", api.GetProductsLog)
	// 輸出 csv
	productsGroup.GET("/export", api.ExportCsv)
	// 取得圖片列表
	productsGroup.GET("/getPictures", api.GetPictureLists)
	// 輸出圖片
	productsGroup.GET("/showPicture/:id", api.ShowPicture)
	// 上傳圖片
	productsGroup.POST("/uploadsPic", api.UploadPic)
	// 更新圖片
	productsGroup.PUT("/updatePic", api.UpdatePic)
	// 刪除圖片
	productsGroup.DELETE("/deletePic/:id", api.DeletePic)
	// 取得產品累吧列表
	productsGroup.GET("/type", api.GetProductType)
	productsGroup.POST("/type", api.AddProductType)
	productsGroup.PUT("/type", api.UpdateProductType)
	productsGroup.DELETE("/type/:id", api.DeleteProductType)
	// 排序商品
	productsGroup.POST("/sort", api.SortProducts)
	// === orders Groups ===
	ordersGroup := apiGroup.Group("/orders")
	// 取得訂購列表
	ordersGroup.GET("/getLists", api.GetOrdersLists)
	// 新增訂購單
	ordersGroup.POST("/add", api.AddOrders)
	// 取得訂單資料
	ordersGroup.GET("/:id", api.GetOrders)
	// 編輯訂購單
	ordersGroup.PATCH("/edit", api.EditOrders)
	// 編輯訂購單
	ordersGroup.DELETE("/:id", api.DeleteOrders)
	// 檢查存貨
	ordersGroup.POST("/check", api.CheckOrders)
	// === setting Groups ===
	settingGroup := apiGroup.Group("/settings")
	settingGroup.GET("/unsetLine/:username", api.UnsetLine)
	// === banner Group ===
	apiGroup.GET("/banner", api.GetBannerList)
	apiGroup.POST("/banner", api.UploadBanner)
	apiGroup.PUT("/banner", api.UpdateBanner)
	apiGroup.DELETE("/banner/:id", api.DeleteBanner)
	// === about Group
	apiGroup.POST("/about", api.UpdateAbout)
	apiGroup.POST("/ckeditor/uploads", api.CkeditorUploads)
	// === contact Group ===
	apiGroup.GET("/contact", api.GetContactList)
	apiGroup.DELETE("/contact/:id", api.DeleteContact)

	// === front ===
	frontGroup := apiGroup.Group("/front")
	frontGroup.GET("/products", front.GetProductsList)
	frontGroup.POST("/contact", middleware.CsrfAuth(), front.InsertContact)
}
