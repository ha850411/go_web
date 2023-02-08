package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	username, _ := c.Get("username")
	c.HTML(http.StatusOK, "index", gin.H{
		"active":   "index",
		"username": username,
	})
}
func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", nil)
}
