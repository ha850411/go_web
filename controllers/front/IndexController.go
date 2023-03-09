package front

import (
	"fmt"
	"goWeb/controllers/api"
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
	Content string   `json:"content"`
}

func Index(c *gin.Context) {
	banner, _, _ := api.GetBannerData("10", "0")
	fmt.Printf("banner: %v\n", banner)
	c.HTML(http.StatusOK, "index", gin.H{
		"active": "index",
		"banner": banner,
	})
}

func Product(c *gin.Context) {
	c.HTML(http.StatusOK, "front_product", gin.H{
		"active": "product",
	})
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
		"active":      "detail",
		"id":          id,
		"data":        rowData,
		"relatedData": relatedData,
	})
}

func ShoppingCart(c *gin.Context) {
	c.HTML(http.StatusOK, "front_cart", gin.H{})
}

func About(c *gin.Context) {
	aboutData := struct {
		Title1   string
		Title2   string
		Content1 string
		Content2 string
		Picture1 string
		Picture2 string
	}{}
	db := database.DbConnect()
	err := db.QueryRow("SELECT IFNULL(title1, ''), IFNULL(content1, ''), IFNULL(picture1, ''), IFNULL(title2, ''), IFNULL(content2, ''), IFNULL(picture2, '') FROM about WHERE id=1").Scan(&aboutData.Title1, &aboutData.Content1, &aboutData.Picture1, &aboutData.Title2, &aboutData.Content2, &aboutData.Picture2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("aboutData: %v\n", aboutData)
	c.HTML(http.StatusOK, "front_about", gin.H{
		"active":    "about",
		"aboutData": aboutData,
	})
}

func Contact(c *gin.Context) {
	c.HTML(http.StatusOK, "front_contact", gin.H{
		"active": "contact",
	})
}

func getProductById(id string) (ProductInfo, error) {
	rowData := ProductInfo{}
	db = database.DbConnect()
	err := db.QueryRow(`SELECT id, name, price, IFNULL(content, '') FROM products 
	WHERE id = ?`, id).Scan(&rowData.Id, &rowData.Name, &rowData.Price, &rowData.Content)
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
