package api

import (
	"fmt"
	"goWeb/conf"
	"goWeb/service"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateAbout(c *gin.Context) {
	title1 := c.PostForm("title1")
	title2 := c.PostForm("title2")
	content1 := c.PostForm("content1")
	content2 := c.PostForm("content2")
	updateData := map[string]string{
		"title1":   title1,
		"title2":   title2,
		"content1": content1,
		"content2": content2,
	}
	formFile1, err := c.FormFile("picture1")
	if err == nil {
		file, _ := formFile1.Open()
		defer file.Close()
		pic1Content, _ := ioutil.ReadAll(file)
		pic1FileName := deleteAndUpdatePic("picture1", pic1Content)
		updateData["picture1"] = pic1FileName
	}
	formFile2, err := c.FormFile("picture2")
	if err == nil {
		file, _ := formFile2.Open()
		defer file.Close()
		pic2Content, _ := ioutil.ReadAll(file)
		pic1FileName := deleteAndUpdatePic("picture2", pic2Content)
		updateData["picture2"] = pic1FileName
	}
	// 若不存在資料則新增一筆
	var count int
	db.QueryRow("SELECT count(*) FROM about").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO about (id) VALUES (1)")
	}

	var setSlice []string
	for k, v := range updateData {
		setSlice = append(setSlice, fmt.Sprintf("%v='%v'", k, v))
	}
	setStr := strings.Join(setSlice, ",")
	sql := fmt.Sprintf("UPDATE about SET %s WHERE id=1", setStr)
	fmt.Printf("sql: %v\n", sql)
	rows, err2 := db.Exec(sql)
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
	}
	RowsAffected, _ := rows.RowsAffected()
	fmt.Printf("RowsAffected: %v\n", RowsAffected)
}

func deleteAndUpdatePic(field string, content []byte) string {
	// 刪除舊檔案
	var picture1Name string
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/about"
	db.QueryRow(`SELECT ? FROM about`, field).Scan(&picture1Name)
	if picture1Name != "" {
		os.Remove(targetDir + "/" + picture1Name)
	}
	// 更新檔案
	var uuid string
	db.QueryRow("SELECT uuid_short()").Scan(&uuid)
	fileName := uuid + ".jpg"
	service.WriteFile(targetDir, fileName, content)
	return fileName
}

func CkeditorUploads(c *gin.Context) {
	formFile, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
	}
	file, _ := formFile.Open()
	defer file.Close()
	picContent, _ := ioutil.ReadAll(file)
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/ckeditor"
	var uuid string
	db.QueryRow("SELECT uuid_short()").Scan(&uuid)
	fileName := uuid + ".jpg"
	service.WriteFile(targetDir, fileName, picContent)
	c.JSON(http.StatusOK, gin.H{
		"url": "/static/uploads/ckeditor/" + fileName,
	})
}
