package front

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{})
}

func Product(c *gin.Context) {
	c.HTML(http.StatusOK, "front_product", gin.H{})
}

func ProductDetail(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "front_product_detail", gin.H{
		id: id,
	})
}

func ShoppingCart(c *gin.Context) {
	c.HTML(http.StatusOK, "front_cart", gin.H{})
}

func About(c *gin.Context) {
	c.HTML(http.StatusOK, "front_about", gin.H{})
}

func Contact(c *gin.Context) {
	c.HTML(http.StatusOK, "front_contact", gin.H{})
}
