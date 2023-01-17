package routes

import (
	"github.com/gin-gonic/gin"
	"goWeb/controllers"
)

func SetRouter(r *gin.Engine) {
	// 載入 web routers
	GetWebRouters(r)
	// 載入 api routers
	GetApiRouters(r)
	// 載入 error routers
	r.NoRoute(controllers.NotFound)
}
