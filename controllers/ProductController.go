package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductManage(c *gin.Context) {
	test := c.MustGet("Test").(string)
	fmt.Printf("test: %v\n", test)
	c.HTML(http.StatusOK, "product", gin.H{
		"active": "product",
		"hello":  test,
	})
}
