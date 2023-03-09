package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetContactList(c *gin.Context) {
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數

	var count int
	sql := "SELECT count(*) FROM contact"
	db.QueryRow(sql).Scan(&count)

	sql = fmt.Sprintf("SELECT id, name, contact, message, createTime FROM contact ORDER BY createTime desc LIMIT %s, %s", offset, limit)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close()
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := struct {
			Id         int       `json:"id"`
			Name       string    `json:"name"`
			Contact    string    `json:"contact"`
			Message    string    `json:"message"`
			CreateTime time.Time `json:"createTime"`
			FormatTime string    `json:"formatTime"`
		}{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Contact, &rowData.Message, &rowData.CreateTime)
		rowData.FormatTime = rowData.CreateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func DeleteContact(c *gin.Context) {
	id := c.Param("id")
	sql := fmt.Sprintf("DELETE FROM contact WHERE id=%s", id)
	_, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err,
		})
		log.Panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "success",
	})
}
