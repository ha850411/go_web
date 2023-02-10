package controllers

import (
	"encoding/base64"
	"fmt"
	"goWeb/database"
	"goWeb/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettingManage(c *gin.Context) {
	username, _ := c.Get("username")
	output := service.GetCommonOutput() // 取得 menu
	output["active"] = "setting"
	output["username"] = username
	c.HTML(http.StatusOK, "setting", output)
}

func UpdatePassword(c *gin.Context) {
	token, _ := c.Cookie("loginToken")
	password := c.PostForm("password")
	username, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	sql := fmt.Sprintf("UPDATE users SET password = '%s' WHERE username='%s'", password, username)
	db := database.DbConnect()
	rows, err2 := db.Exec(sql)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 0 {
		c.SetCookie("loginToken", "", -1, "/", "", false, false) // 登出
		c.Redirect(http.StatusMovedPermanently, "/admin/login")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "變更密碼失敗",
		})
	}
}
