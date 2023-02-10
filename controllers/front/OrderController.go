package front

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	c.HTML(http.StatusOK, "front-order", gin.H{})
}
