package front

import (
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
		"pid":     c.PostForm("pid"),
		"amount":  c.PostForm("amount"),
		"remark":  c.PostForm("remark"),
	}
	sql := fmt.Sprintf("INSERT INTO orders (name, contact, pid, amount, remark) VALUES ('%s', '%s', %s, %s, '%s')", postData["name"], postData["contact"], postData["pid"], postData["amount"], postData["remark"])
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
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<script>alert('新增成功');location='/orders'</script>`)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}
