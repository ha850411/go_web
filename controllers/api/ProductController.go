package api

import (
	"fmt"
	"goWeb/database"
	"goWeb/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProductsList(c *gin.Context) {
	db := database.DbConnect()
	// params
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數
	keyword := c.Query("search[value]")     // 關鍵字
	where := "1"
	if keyword != "" {
		where = fmt.Sprintf("name LIKE '%%%s%%'", keyword)
	}
	// 定義排序欄位
	columes := map[string]string{
		"0": "name",
		"1": "amount",
		"2": "updateTime",
	}
	orderBy := c.DefaultQuery("order[0][column]", "0")   // 分頁筆數
	orderType := c.DefaultQuery("order[0][dir]", "desc") // 起始筆數
	// count total
	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM products WHERE %s", where)
	db.QueryRow(sql).Scan(&count)
	// 分頁
	sql = fmt.Sprintf("SELECT * FROM products WHERE %s ORDER BY %s %s LIMIT %s OFFSET %s", where, columes[orderBy], orderType, limit, offset)
	rows, err := db.Query(sql)
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
		rowData := models.Products{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Amount, &rowData.UpdateTime)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}
