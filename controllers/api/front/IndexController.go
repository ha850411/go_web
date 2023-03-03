package front

import (
	"database/sql"
	"goWeb/database"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func init() {
	db = database.DbConnect()
}
func GetProductsList(c *gin.Context) {
	var count int
	page := c.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)
	perpage := 9
	db.QueryRow(`SELECT count(*) FROM products`).Scan(&count)
	// data
	rows, err := db.Query(`SELECT id, name, price FROM products 
	WHERE 1 ORDER BY id desc LIMIT ? OFFSET ?`, perpage, perpage*pageInt-perpage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := struct {
			Id      int      `json:"id"`
			Name    string   `json:"name"`
			Price   int      `json:"price"`
			Picture []string `json:"picture"`
		}{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Price)
		var pictures string
		db.QueryRow(`SELECT GROUP_CONCAT(DISTINCT picture) as pictures FROM products_picture WHERE pid= ?`, rowData.Id).Scan(&pictures)
		if pictures != "" {
			rowData.Picture = strings.Split(pictures, ",")
		}
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"data":  data,
	})
}
