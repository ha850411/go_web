package front

import (
	"encoding/json"
	"fmt"
	"goWeb/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	c.HTML(http.StatusOK, "front-order", gin.H{})
}

func OrdersAdd(c *gin.Context) {
	postData := map[string]string{
		"name":    c.PostForm("name"),
		"contact": c.PostForm("contact"),
		"remark":  c.PostForm("remark"),
		"detail":  c.PostForm("detail"),
	}
	sql := fmt.Sprintf("INSERT INTO orders (name, contact, remark) VALUES ('%s', '%s', '%s')", postData["name"], postData["contact"], postData["remark"])
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
		for _, jsonMap := range detail {
			sql := fmt.Sprintf("INSERT INTO orders_detail (order_id, pid, amount) VALUES (%v, %v, %v)", LastInsertId, jsonMap["pid"], jsonMap["amount"])
			db.Exec(sql)
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<script>alert('新增成功');location='/orders'</script>`)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}
