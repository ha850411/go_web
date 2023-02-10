package api

import (
	"fmt"
	"goWeb/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrdersLists(c *gin.Context) {
	// params
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數
	keyword := c.Query("search[value]")     // 關鍵字
	where := "1"
	if keyword != "" {
		where = fmt.Sprintf("a.name LIKE '%%%s%%' OR a.contact LIKE '%%%s%%' OR a.remark LIKE '%%%s%%' OR b.name LIKE '%%%s%%'", keyword, keyword, keyword, keyword)
	}
	// 定義排序欄位
	columes := map[string]string{
		"0": "a.id",
		"1": "a.name",
		"2": "a.contact",
		"3": "a.pid",
		"4": "a.amount",
		"5": "a.remark",
		"6": "a.updateTime",
	}
	orderBy := c.DefaultQuery("order[0][column]", "0")   // 分頁筆數
	orderType := c.DefaultQuery("order[0][dir]", "desc") // 起始筆數
	// count total
	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM orders as a LEFT JOIN products as b ON a.pid=b.id WHERE %s", where)
	fmt.Printf("sql: %v\n", sql)
	db.QueryRow(sql).Scan(&count)
	// 分頁
	sql = fmt.Sprintf("SELECT a.*, b.name as pname FROM orders as a LEFT JOIN products as b ON a.pid=b.id WHERE %s ORDER BY %s %s LIMIT %s OFFSET %s", where, columes[orderBy], orderType, limit, offset)
	fmt.Printf("sql: %v\n", sql)
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
		rowData := models.Orders{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Contact, &rowData.Pid, &rowData.Amount, &rowData.Remark, &rowData.UpdateTime, &rowData.Pname)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		fmt.Printf("rowData: %v\n", rowData)
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func AddOrders(c *gin.Context) {
	name := c.PostForm("name")
	contact := c.PostForm("contact")
	pid := c.PostForm("pid")
	amount := c.PostForm("amount")
	remark := c.DefaultPostForm("remark", "")
	sql := fmt.Sprintf("INSERT INTO orders (name, contact, pid, amount, remark) VALUES ('%s', '%s', %s, %s, '%s')", name, contact, pid, amount, remark)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := rows.RowsAffected()
	if rowAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":         http.StatusOK,
			"rowsAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "新增失敗",
		})
	}
}

func GetOrders(c *gin.Context) {
	id := c.Param("id")
	sql := fmt.Sprintf("SELECT a.*, b.name as pname FROM orders as a LEFT JOIN products as b ON a.pid=b.id WHERE a.id=%s", id)

	rowData := models.Orders{}
	err := db.QueryRow(sql).Scan(&rowData.Id, &rowData.Name, &rowData.Contact, &rowData.Pid, &rowData.Amount, &rowData.Remark, &rowData.UpdateTime, &rowData.Pname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": rowData,
	})
}

func EditOrders(c *gin.Context) {
	postData := map[string]string{
		"id":      c.PostForm("id"),
		"name":    c.PostForm("name"),
		"contact": c.PostForm("contact"),
		"pid":     c.PostForm("pid"),
		"amount":  c.PostForm("amount"),
		"remark":  c.PostForm("remark"),
	}
	sql := fmt.Sprintf("UPDATE orders SET name='%s', contact='%s', pid=%s, amount=%s, remark='%s' WHERE id=%s", postData["name"], postData["contact"], postData["pid"], postData["amount"], postData["remark"], postData["id"])
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":         http.StatusOK,
			"rowsAffected": rowsAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func DeleteOrders(c *gin.Context) {
	id := c.Param("id")
	sql := fmt.Sprintf("DELETE FROM orders WHERE id=%s", id)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":         http.StatusOK,
			"rowsAffected": rowsAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}
