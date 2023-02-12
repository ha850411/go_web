package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Analysis(c *gin.Context) {

	output := service.GetCommonOutput(c, "analysis") // 取得 menu
	c.HTML(http.StatusOK, "analysis", output)
}
