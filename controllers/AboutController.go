package controllers

import (
	"fmt"
	"goWeb/database"
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AboutManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "about") // 取得 menu
	aboutData := struct {
		Title1   string
		Title2   string
		Content1 string
		Content2 string
		Picture1 string
		Picture2 string
	}{}
	db := database.DbConnect()
	err := db.QueryRow("SELECT IFNULL(title1, ''), IFNULL(content1, ''), IFNULL(picture1, ''), IFNULL(title2, ''), IFNULL(content2, ''), IFNULL(picture2, '') FROM about WHERE id=1").Scan(&aboutData.Title1, &aboutData.Content1, &aboutData.Picture1, &aboutData.Title2, &aboutData.Content2, &aboutData.Picture2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("aboutData: %v\n", aboutData)
	output["aboutData"] = aboutData
	c.HTML(http.StatusOK, "aboutManage", output)
}
