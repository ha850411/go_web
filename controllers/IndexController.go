package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", nil)
}
