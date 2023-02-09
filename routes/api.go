package routes

import (
	"goWeb/controllers/api"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	{
		// 取得商品列表
		apiGroup.GET("/products/getList", api.GetProductsList)
		// 更新數量
		apiGroup.PUT("/products/updateAmount", api.UpdateAmount)
		// 新增產品
		apiGroup.POST("/products/add", api.AddProduct)
		// 編輯產品
		apiGroup.PATCH("/products/edit", api.EditProduct)
		// 刪除產品
		apiGroup.DELETE("/products/delete/:id", api.DeleteProduct)
		// 取得存貨警戒商品
		apiGroup.GET("/products/getTips", api.GetTips)
		// 取得商品名稱
		apiGroup.GET("/products/getProductsNameList", api.GetProductsNameList)
		// 取得商品 log
		apiGroup.GET("/products/getProductsLog/:pid", api.GetProductsLog)

	}
}
