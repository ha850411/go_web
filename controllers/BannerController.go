package controllers

import (
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BannerManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "banner") // 取得 menu
	c.HTML(http.StatusOK, "bannerManage", output)
}
