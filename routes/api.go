package routes

import (
	"goWeb/controllers/api"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	{
		// 取得商品列表
		apiGroup.GET("/getProductsList", api.GetProductsList)
	}
}
