package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductManage(c *gin.Context) {
	c.HTML(http.StatusOK, "product", gin.H{
		"active": "product",
	})
}
