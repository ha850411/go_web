package routes

import (
	"github.com/gin-gonic/gin"
	controllerApi "goWeb/controllers/api"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	{
		apiGroup.GET("/test", controllerApi.Test)
	}
}
