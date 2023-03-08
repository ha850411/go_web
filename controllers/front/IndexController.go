package front

import (
	"fmt"
	"goWeb/controllers/api/front"
	"goWeb/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductInfo struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Price   int      `json:"price"`
	Picture []string `json:"picture"`
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{})
}

func Product(c *gin.Context) {
	c.HTML(http.StatusOK, "front_product", gin.H{})
}

func ProductDetail(c *gin.Context) {
	id := c.Param("id")
	// 取得商品資訊
	rowData, err := getProductById(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// 取得推薦商品
	relatedData, _ := front.GetFrontProducts(1, "")
	c.HTML(http.StatusOK, "front_product_detail", gin.H{
		"id":          id,
		"data":        rowData,
		"relatedData": relatedData,
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

func getProductById(id string) (ProductInfo, error) {
	rowData := ProductInfo{}
	db = database.DbConnect()
	err := db.QueryRow(`SELECT id, name, price FROM products 
	WHERE id = ?`, id).Scan(&rowData.Id, &rowData.Name, &rowData.Price)
	if err != nil {
		return rowData, err
	}
	var pictures string
	db.QueryRow(`SELECT GROUP_CONCAT(DISTINCT picture) as pictures FROM products_picture WHERE pid= ?`, id).Scan(&pictures)
	fmt.Printf("pictures: %v\n", pictures)
	if pictures != "" {
		rowData.Picture = strings.Split(pictures, ",")
		for i, v := range rowData.Picture {
			rowData.Picture[i] = fmt.Sprintf("/static/uploads/%v/%v", id, v)
		}
	}
	fmt.Printf("rowData: %v\n", rowData)
	return rowData, nil
}
