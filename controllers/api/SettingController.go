package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnsetLine(c *gin.Context) {
	username := c.Param("username")
	sql := fmt.Sprintf("UPDATE users SET linebot_token=NULL WHERE username='%s'", username)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	if rowAffected, _ := row.RowsAffected(); rowAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "重置失敗",
		})
	}
}
