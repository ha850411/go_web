package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"active": "index",
	})
}
func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", nil)
}
