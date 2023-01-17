package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"view": "test",
	})
}
func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}
