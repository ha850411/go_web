package routes

import (
	"goWeb/controllers"

	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	// 載入 web routers
	GetWebRouters(r)
	// 載入 api routers
	GetApiRouters(r)
	// 載入 error routers
	r.NoRoute(controllers.NotFound)
	// 加載資源文件
	r.Static("/static", "./static")
}
