package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ContactManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "contact") // 取得 menu
	c.HTML(http.StatusOK, "contactManage", output)
}
