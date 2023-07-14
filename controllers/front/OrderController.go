package front

import (
	"encoding/json"
	"fmt"
	"goWeb/database"
	"goWeb/service/linebot"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	c.HTML(http.StatusOK, "front-order", gin.H{})
}

func OrdersAdd(c *gin.Context) {
	type PostData struct {
		Name    string `json:"name"`
		Contact string `json:"contact"`
		Address string `json:"address"`
		Remark  string `json:"remark"`
		Detail  string `json:"detail"`
	}
	postData := PostData{}
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &postData)
	fmt.Printf("postData: %v\n", postData)

	sql := fmt.Sprintf("INSERT INTO orders (name, contact, address, remark) VALUES ('%s', '%s','%s', '%s')", postData.Name, postData.Contact, postData.Address, postData.Remark)
	fmt.Printf("sql: %v\n", sql)
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
		_ = json.Unmarshal([]byte(postData.Detail), &detail)
		LastInsertId, _ := rows.LastInsertId()
		lineMessage := fmt.Sprintf("\n*新訂單\n姓名: %s\n聯絡方式: %s\n配送地址: %s\n備註: %s\n訂購商品:", postData.Name, postData.Contact, postData.Address, postData.Remark)
		var total int
		for _, jsonMap := range detail {
			// get price
			var price, originPrice, discountPrice int

			amount := int(jsonMap["amount"].(float64))
			db.QueryRow(`SELECT price, discount_price FROM products WHERE id= ?`, jsonMap["id"]).Scan(&originPrice, &discountPrice)
			price = originPrice
			if discountPrice > 0 {
				price = discountPrice
			}

			// do insert
			sql := fmt.Sprintf("INSERT INTO orders_detail (order_id, pid, amount, price, total) VALUES (%v, %v, %v, %v, %v)", LastInsertId, jsonMap["id"], jsonMap["amount"], price, price*amount)
			db.Exec(sql)
			// 推播訊息
			sql = fmt.Sprintf("SELECT name FROM products WHERE id = %v", jsonMap["id"])
			var pname string
			db.QueryRow(sql).Scan(&pname)
			lineMessage += fmt.Sprintf("\n  商品: %s\n  數量: %v\n  單價: %v\n", pname, amount, price)
			// 計算總價
			total += price * amount
		}
		lineMessage += fmt.Sprintf("===============\n總價: %v", total)
		// 推播
		linebot.Request(lineMessage)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "新增成功",
		})
		// c.Header("Content-Type", "text/html; charset=utf-8")
		// c.String(200, `<script>alert('訂購成功');location='/orders'</script>`)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "新增失敗",
		})
	}
}
