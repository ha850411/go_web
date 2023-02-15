package controllers

import (
	"fmt"
	"goWeb/database"
	"goWeb/models"
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product") // 取得 menu
	keyword := c.DefaultQuery("keyword", "")
	output["keyword"] = keyword
	productsType := getProductType()
	fmt.Printf("productsType: %+v\n", productsType)
	output["productsType"] = productsType
	c.HTML(http.StatusOK, "product", output)
}

func ProductPicManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product") // 取得 menu
	pid := c.Param("pid")
	output["pid"] = pid
	output["pname"] = getProductName(pid)
	fmt.Printf("pid: %v\n", pid)
	c.HTML(http.StatusOK, "product.picture", output)
}

func ProductTypeManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product.type") // 取得 menu
	c.HTML(http.StatusOK, "product.type", output)
}

/*
* 取商品名稱 by pid
 */
func getProductName(id string) (pname string) {
	db := database.DbConnect()
	db.QueryRow(`SELECT name FROM products WHERE id=?`, id).Scan(&pname)
	return
}

/*
* 取得商品類別 list
 */
func getProductType() []interface{} {
	db := database.DbConnect()
	rows, _ := db.Query(`SELECT id, name FROM products_type`)
	data := make([]interface{}, 0)
	defer rows.Close()
	for rows.Next() {
		var temp models.ProductsType
		rows.Scan(&temp.Id, &temp.Name)
		data = append(data, temp)
	}
	return data
}
