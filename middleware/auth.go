package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 中間件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("loginToken")
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
		c.Next()
	}
}
