package api

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdatePassword(c *gin.Context) {
	token, _ := c.Cookie("loginToken")
	username, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	password := c.PostForm("password")
	sql := fmt.Sprintf("UPDATE users SET password = '%s' WHERE username='%s'", password, username)
	rows, err2 := db.Exec(sql)
	if err2 != nil {
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
		c.SetCookie("loginToken", "", 0, "/", "", false, false) // 登出
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "變更密碼失敗",
		})
	}
}
