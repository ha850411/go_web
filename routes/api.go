package routes

import (
	controllerApi "goWeb/controllers/api"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(r *gin.Engine) {
	apiGroup := r.Group("api")
	{
		apiGroup.GET("/test", controllerApi.Test)
	}
}
