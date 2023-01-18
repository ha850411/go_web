package routes

import (
	"goWeb/controllers"

	"github.com/gin-gonic/gin"
)

func GetWebRouters(r *gin.Engine) {
	r.GET("/", controllers.Index)
	r.GET("/index", controllers.Index)
	r.GET("/product", controllers.ProductManage)
}
