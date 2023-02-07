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
	}
}
