package api

import (
	"database/sql"
	"fmt"
	"goWeb/database"
	"goWeb/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func init() {
	db = database.DbConnect()
}

func GetProductsList(c *gin.Context) {
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
		"2": "amountNotice",
		"3": "updateTime",
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
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Amount, &rowData.AmountNotice, &rowData.UpdateTime)
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

func UpdateAmount(c *gin.Context) {
	id := c.PostForm("id")
	amount := c.PostForm("amount")
	sql := fmt.Sprintf("UPDATE products SET amount=%s WHERE id=%s", amount, id)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowsAffected, _ := rows.RowsAffected()
	fmt.Printf("sql: %v\n", sql)
	fmt.Printf("rowsAffected: %v\n", rowsAffected)
	if rowsAffected > 0 {
		writeProductsLog(id, amount) // 寫 log
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"updatedRows": rowsAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func AddProduct(c *gin.Context) {
	name := c.PostForm("name")
	amount := c.PostForm("amount")
	amountNotice := c.DefaultPostForm("amount_notice", "0")
	if len(amountNotice) == 0 {
		amountNotice = "0"
	}
	fmt.Printf("name: %v\n", name)
	fmt.Printf("amount: %v\n", amount)
	fmt.Printf("amountNotice: %v\n", amountNotice)
	sql := fmt.Sprintf("INSERT INTO products (name, amount, amountNotice) VALUES ('%s', %s, %s)", name, amount, amountNotice)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	lastInsertID, insertErr := row.LastInsertId()
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  insertErr,
		})
		log.Panic(insertErr)
	}
	writeProductsLog(strconv.FormatInt(lastInsertID, 10), amount)
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"insertId": lastInsertID,
	})
}

func EditProduct(c *gin.Context) {
	editId := c.PostForm("id")
	name := c.PostForm("name")
	amount := c.PostForm("amount")
	amountNotice := c.DefaultPostForm("amount_notice", "0")
	if len(amountNotice) == 0 {
		amountNotice = "0"
	}
	// 檢查數量是否有變
	var nowAmount string
	sql := fmt.Sprintf("SELECT amount FROM products WHERE id=%s", editId)
	db.QueryRow(sql).Scan(&nowAmount)
	fmt.Printf("amount: %v\n", amount)
	fmt.Printf("nowAmount: %v\n", nowAmount)

	sql = fmt.Sprintf("UPDATE products SET name='%s', amount=%s, amountNotice=%s WHERE id=%s", name, amount, amountNotice, editId)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := row.RowsAffected()
	if rowAffected > 0 {
		if amount != nowAmount {
			writeProductsLog(editId, amount)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("id: %v\n", id)
	sql := fmt.Sprintf("DELETE FROM products WHERE id=%s", id)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := row.RowsAffected()
	if rowAffected > 0 {
		// 刪除 log
		sql = fmt.Sprintf("DELETE FROM products_log WHERE pid=%s", id)
		db.Exec(sql)
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}

// 寫 log
func writeProductsLog(id string, amount string) {
	today := time.Now().Format("2006-01-02")
	var check int
	sql := fmt.Sprintf("SELECT count(*) FROM products_log WHERE pid=%s AND updateDate='%s'", id, today)
	db.QueryRow(sql).Scan(&check)
	if check > 0 {
		sql = fmt.Sprintf("UPDATE products_log SET amount=%s WHERE pid=%s AND updateDate='%s'", amount, id, today)
		db.Exec(sql)
	} else {
		sql = fmt.Sprintf("INSERT INTO products_log (pid, amount, updateDate) VALUES (%s, %s, '%s')", id, amount, today)
		db.Exec(sql)
	}
}

func GetTips(c *gin.Context) {
	sql := "SELECT * FROM products WHERE amount<=amountNotice"
	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := models.Products{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Amount, &rowData.AmountNotice, &rowData.UpdateTime)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
