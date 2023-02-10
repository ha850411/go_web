package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderManage(c *gin.Context) {
	username, _ := c.Get("username")
	output := service.GetCommonOutput() // 取得 menu
	output["active"] = "order"
	output["username"] = username
	c.HTML(http.StatusOK, "order", output)
}
