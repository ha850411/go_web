package front

import (
	"database/sql"
	"fmt"
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

// - 取得商品列表
func GetProductsList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)
	keyword := c.DefaultQuery("keyword", "")
	kind := c.DefaultQuery("kind", "all")
	rand := false
	if c.Query("rand") == "Y" {
		rand = true
	}
	data, count := GetFrontProducts(pageInt, keyword, kind, rand)
	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"data":  data,
	})
}

func GetFrontProducts(page int, keyword string, kind string, orderbyRand bool) ([]interface{}, int) {
	var count int
	perpage := 8
	where := ""
	if keyword != "" {
		where = fmt.Sprintf(" AND name LIKE '%%%s%%'", keyword)
	}
	if kind != "all" && kind != "" {
		where += fmt.Sprintf(" AND type=%v", kind)
	}
	orderby := "name asc"
	if orderbyRand {
		orderby = "RAND()"
	}
	fmt.Printf("keyword: %v\n", keyword)
	fmt.Printf("where: %v\n", where)
	sql := fmt.Sprintf("SELECT count(*) FROM products WHERE status=1 %s", where)
	db.QueryRow(sql).Scan(&count)
	// data
	sql = fmt.Sprintf("SELECT id, name, price FROM products WHERE status=1 %v ORDER BY %s LIMIT %v OFFSET %v", where, orderby, perpage, perpage*page-perpage)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
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
	return data, count
}

func InsertContact(c *gin.Context) {
	var param struct {
		Name    string `json:"name"`
		Contact string `json:"contact"`
		Message string `json:"message"`
	}
	c.Bind(&param)
	sql := fmt.Sprintf("INSERT INTO contact (name, contact, message) VALUES ('%v', '%v', '%v')", param.Name, param.Contact, param.Message)
	fmt.Printf("sql: %v\n", sql)
	_, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":         200,
		"AffectedRows": 1,
	})
}
