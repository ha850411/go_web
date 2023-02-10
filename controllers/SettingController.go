package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettingManage(c *gin.Context) {
	username, _ := c.Get("username")
	output := service.GetCommonOutput() // 取得 menu
	output["active"] = "setting"
	output["username"] = username
	c.HTML(http.StatusOK, "setting", output)
}
