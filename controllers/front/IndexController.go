package front

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{})
}

func Product(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "product", gin.H{
		id: id,
	})
}
