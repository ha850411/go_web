package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "order") // 取得 menu
	c.HTML(http.StatusOK, "order", output)
}
