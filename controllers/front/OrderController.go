package front

import (
	"encoding/json"
	"fmt"
	"goWeb/database"
	"goWeb/service/linebot"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	c.HTML(http.StatusOK, "front-order", gin.H{})
}

func OrdersAdd(c *gin.Context) {
	postData := map[string]string{
		"name":    c.PostForm("name"),
		"contact": c.PostForm("contact"),
		"address": c.PostForm("address"),
		"remark":  c.PostForm("remark"),
		"detail":  c.PostForm("detail"),
	}
	sql := fmt.Sprintf("INSERT INTO orders (name, contact, address, remark) VALUES ('%s', '%s','%s', '%s')", postData["name"], postData["contact"], postData["address"], postData["remark"])
	db := database.DbConnect()
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 0 {
		// 解析 json 並新增訂單內容
		var detail map[string]map[string]interface{}
		_ = json.Unmarshal([]byte(postData["detail"]), &detail)
		LastInsertId, _ := rows.LastInsertId()
		lineMessage := fmt.Sprintf("\n*新訂單\n姓名: %s\n聯絡方式: %s\n配送地址: %s\n備註: %s\n訂購商品:", postData["name"], postData["contact"], postData["address"], postData["remark"])
		for _, jsonMap := range detail {
			// get price
			var price int
			amount, _ := strconv.Atoi(jsonMap["amount"].(string))
			db.QueryRow(`SELECT price FROM products WHERE id= ?`, jsonMap["pid"]).Scan(&price)
			// do insert
			sql := fmt.Sprintf("INSERT INTO orders_detail (order_id, pid, amount, price, total) VALUES (%v, %v, %v, %v, %v)", LastInsertId, jsonMap["pid"], jsonMap["amount"], price, price*amount)
			db.Exec(sql)
			// 推播訊息
			sql = fmt.Sprintf("SELECT name FROM products WHERE id = %s", jsonMap["pid"])
			var pname string
			db.QueryRow(sql).Scan(&pname)
			lineMessage += fmt.Sprintf("\n- 商品: %s\n  數量: %s", pname, jsonMap["amount"])
		}
		// 推播
		linebot.Request(lineMessage)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<script>alert('訂購成功');location='/orders'</script>`)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "新增失敗",
		})
	}
}
