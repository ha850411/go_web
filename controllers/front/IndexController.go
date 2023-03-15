package front

import (
	"fmt"
	"goWeb/controllers/api"
	"goWeb/controllers/api/front"
	"goWeb/database"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductInfo struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Price   int      `json:"price"`
	Picture []string `json:"picture"`
	Content string   `json:"content"`
	Type    int      `json:"type"`
}

func GetCommonOutput(active string) map[string]interface{} {
	output := make(map[string]interface{})
	output["active"] = active
	output["staticFreshFlag"] = time.Now().Unix()
	return output
}

func Index(c *gin.Context) {
	banner, _, _ := api.GetBannerData("10", "0")
	output := GetCommonOutput("index")
	output["banner"] = banner
	c.HTML(http.StatusOK, "index", output)
}

func Product(c *gin.Context) {
	output := GetCommonOutput("product")
	kind := c.DefaultQuery("kind", "all")
	var categoryTitle string
	switch kind {
	case "wine":
		categoryTitle = "酒類"
	default:
		categoryTitle = "所有商品"
	}
	output["kind"] = kind
	output["categoryTitle"] = categoryTitle
	c.HTML(http.StatusOK, "front_product", output)
}

func ProductDetail(c *gin.Context) {
	id := c.Param("id")
	// 取得商品資訊
	rowData, err := getProductById(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// 取得推薦商品
	relatedData, _ := front.GetFrontProducts(1, "", "all", true)
	output := GetCommonOutput("detail")
	output["id"] = id
	output["data"] = rowData
	output["relatedData"] = relatedData
	c.HTML(http.StatusOK, "front_product_detail", output)
}

func ShoppingCart(c *gin.Context) {
	output := GetCommonOutput("cart")
	c.HTML(http.StatusOK, "front_cart", output)
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
	output := GetCommonOutput("contact")
	c.HTML(http.StatusOK, "front_contact", output)
}

func getProductById(id string) (ProductInfo, error) {
	rowData := ProductInfo{}
	db = database.DbConnect()
	err := db.QueryRow(`SELECT id, name, price, IFNULL(content, ''), type FROM products 
	WHERE id = ?`, id).Scan(&rowData.Id, &rowData.Name, &rowData.Price, &rowData.Content, &rowData.Type)
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
