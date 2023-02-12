package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product") // 取得 menu
	keyword := c.DefaultQuery("keyword", "")
	output["keyword"] = keyword
	c.HTML(http.StatusOK, "product", output)
}
