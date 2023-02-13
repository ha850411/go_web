package controllers

import (
	"encoding/base64"
	"fmt"
	"goWeb/database"
	"goWeb/service"
	"goWeb/service/linebot"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettingManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "setting") // 取得 menu
	var lineUrl string
	// 確認是否綁定 line notify
	var linebotToken string
	sql := fmt.Sprintf("SELECT linebot_token FROM users WHERE username='%s'", output["username"])
	db := database.DbConnect()
	db.QueryRow(sql).Scan(&linebotToken)
	isBind := false
	if linebotToken != "" {
		isBind = true
	}
	lineUrl = fmt.Sprintf("https://notify-bot.line.me/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=notify&state=csrfToken", linebot.CLIENT_ID, linebot.REDIRECT_URI)
	output["lineUrl"] = lineUrl
	output["isBind"] = isBind
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
