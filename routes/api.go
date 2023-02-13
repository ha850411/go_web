package routes

import (
	"goWeb/controllers/api"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	// === products Groups ===
	productsGroup := apiGroup.Group("/products")
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
	// === setting Groups ===
	settingGroup := apiGroup.Group("/settings")
	settingGroup.GET("/unsetLine/:username", api.UnsetLine)
}
