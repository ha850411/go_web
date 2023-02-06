package controllers

import (
	"encoding/base64"
	"fmt"
	"goWeb/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}

func Auth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// db connect
	db := database.DbConnect()
	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM users WHERE username='%s' AND password='%s'", username, password)
	err := db.QueryRow(sql).Scan(&count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
	}
	// base64
	token := base64.RawURLEncoding.EncodeToString([]byte(username))
	if count >= 1 {
		// 記住我
		remember := c.PostForm("remember")
		expire := 3600
		if remember == "Y" {
			expire = 86400 * 30
		}
		c.SetCookie("loginToken", token, expire, "/", "", false, false)
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "帳號密碼錯誤, 請重新登入",
		})
	}
}
