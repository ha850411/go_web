package api

import (
	"encoding/json"
	"fmt"
	"goWeb/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetOrdersLists(c *gin.Context) {
	// params
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數
	keyword := c.Query("search[value]")     // 關鍵字
	where := "1"
	if keyword != "" {
		where = fmt.Sprintf("a.name LIKE '%%%s%%' OR a.contact LIKE '%%%s%%' OR a.remark LIKE '%%%s%%' OR c.name LIKE '%%%s%%'", keyword, keyword, keyword, keyword)
	}
	// 定義排序欄位
	columes := map[string]string{
		"0": "name",
		"1": "contact",
		"3": "remark",
		"4": "updateTime",
	}
	orderBy := c.DefaultQuery("order[0][column]", "0")   // 分頁筆數
	orderType := c.DefaultQuery("order[0][dir]", "desc") // 起始筆數
	// count total
	var ids string
	sql := fmt.Sprintf(`SELECT GROUP_CONCAT( DISTINCT a.id) as ids
	FROM orders as a 
	LEFT JOIN orders_detail as b ON a.id=b.order_id
	LEFT JOIN products as c ON b.pid=c.id
	WHERE %s ORDER BY a.%s %s LIMIT %s OFFSET %s`, where, columes[orderBy], orderType, limit, offset)
	fmt.Printf("sql: %v\n", sql)
	db.QueryRow(sql).Scan(&ids)
	count := len(strings.Split(ids, ","))
	fmt.Printf("ids: %v\n", ids)
	fmt.Printf("count: %v\n", count)
	// 分頁
	sql = fmt.Sprintf("SELECT * FROM orders WHERE id IN (%s) ORDER BY %s %s LIMIT %s OFFSET %s", ids, columes[orderBy], orderType, limit, offset)
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
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Contact, &rowData.Address, &rowData.Remark, &rowData.UpdateTime)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		// 取得 orders_detail - begin
		sql := fmt.Sprintf(`SELECT a.pid, a.amount, b.name ,a.price, a.total
		FROM orders_detail as a
		LEFT JOIN products as b ON a.pid=b.id
		WHERE a.order_id=%v`, rowData.Id)
		fmt.Printf("sql: %v\n", sql)
		detailRows, err2 := db.Query(sql)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "取得商品資訊失敗:" + err2.Error(),
			})
			log.Panic(err)
		}
		defer detailRows.Close() // close
		var total int
		for detailRows.Next() {
			tempDetail := struct {
				Pid    int    `json:"pid"`
				Amount int    `json:"amount"`
				Pname  string `json:"pname"`
				Price  int    `json:"price"`
				Total  int    `json:"total"`
			}{}
			detailRows.Scan(&tempDetail.Pid, &tempDetail.Amount, &tempDetail.Pname, &tempDetail.Price, &tempDetail.Total)
			total += tempDetail.Total
			rowData.Detail = append(rowData.Detail, tempDetail)
		}
		rowData.Total = total
		// 取得 orders_detail - end
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
	address := c.PostForm("address")
	remark := c.DefaultPostForm("remark", "")
	sql := fmt.Sprintf("INSERT INTO orders (name, contact, address, remark) VALUES ('%s', '%s', '%s', '%s')", name, contact, address, remark)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}

	if rowAffected, _ := rows.RowsAffected(); rowAffected > 0 {
		// 解析 json 並新增訂單內容
		tempDetail := c.PostForm("detail")
		var detail map[string]map[string]interface{}
		_ = json.Unmarshal([]byte(tempDetail), &detail)
		LastInsertId, _ := rows.LastInsertId()
		for _, jsonMap := range detail {
			var price int
			amount, _ := strconv.Atoi(jsonMap["amount"].(string))
			db.QueryRow(`SELECT price FROM products WHERE id= ?`, jsonMap["pid"]).Scan(&price)
			sql := fmt.Sprintf("INSERT INTO orders_detail (order_id, pid, amount, price, total) VALUES (%v, %v, %v, %v, %v)", LastInsertId, jsonMap["pid"], jsonMap["amount"], price, price*amount)
			fmt.Printf("sql: %v\n", sql)
			db.Exec(sql)
		}
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
	sql := fmt.Sprintf("SELECT * FROM orders WHERE id=%s", id)

	rowData := models.Orders{}
	err := db.QueryRow(sql).Scan(&rowData.Id, &rowData.Name, &rowData.Contact, &rowData.Address, &rowData.Remark, &rowData.UpdateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	// 取得 orders_detail - begin
	sql = fmt.Sprintf(`SELECT pid, amount, price FROM orders_detail  WHERE order_id=%v`, rowData.Id)
	detailRows, err2 := db.Query(sql)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "取得商品資訊失敗:" + err2.Error(),
		})
		log.Panic(err)
	}
	defer detailRows.Close() // close
	for detailRows.Next() {
		tempDetail := struct {
			Pid    int `json:"pid"`
			Amount int `json:"amount"`
			Price  int `json:"price"`
		}{}
		detailRows.Scan(&tempDetail.Pid, &tempDetail.Amount, &tempDetail.Price)
		rowData.Detail = append(rowData.Detail, tempDetail)
	}
	// 取得 orders_detail - end
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
		"address": c.PostForm("address"),
		"remark":  c.PostForm("remark"),
		"detail":  c.PostForm("detail"),
	}
	sql := fmt.Sprintf("UPDATE orders SET name='%s', contact='%s', address='%s', remark='%s' WHERE id=%s", postData["name"], postData["contact"], postData["address"], postData["remark"], postData["id"])
	fmt.Printf("sql: %v\n", sql)
	_, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}

	// 刪除原始訂單資料
	sql = fmt.Sprintf("DELETE FROM orders_detail WHERE order_id=%v", postData["id"])
	db.Exec(sql)
	// 解析 json 並新增訂單內容
	var detail map[string]map[string]interface{}
	_ = json.Unmarshal([]byte(postData["detail"]), &detail)
	for _, jsonMap := range detail {
		// 取價錢
		var price int
		amount, _ := strconv.Atoi(jsonMap["amount"].(string))
		db.QueryRow(`SELECT price FROM products WHERE id= ?`, jsonMap["pid"]).Scan(&price)
		// do insert
		sql := fmt.Sprintf("INSERT INTO orders_detail (order_id, pid, amount, price, total) VALUES (%v, %v, %v, %v, %v)", postData["id"], jsonMap["pid"], jsonMap["amount"], price, price*amount)
		db.Exec(sql)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":         http.StatusOK,
		"rowsAffected": 1,
	})
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
		sql = fmt.Sprintf("DELETE FROM orders_detail WHERE order_id=%v", id)
		db.Exec(sql)
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
