package api

import (
	"fmt"
	"goWeb/conf"
	"goWeb/models"
	"goWeb/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetBannerList(c *gin.Context) {
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數
	data, count := GetBannerData(limit, offset)
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func GetBannerData(limit string, offset string) ([]interface{}, error) {
	var count int
	db.QueryRow("SELECT count(*) FROM banner").Scan(&count)
	data := make([]interface{}, 0)
	sql := fmt.Sprintf("SELECT id, picture, updateTime FROM banner ORDER BY sort asc LIMIT %s, %s", offset, limit)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		rowData := models.Banner{}
		rows.Scan(&rowData.Id, &rowData.Picture, &rowData.UpdateTime)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	return data, nil
}

func UploadBanner(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	var sort int
	db.QueryRow("SELECT IFNULL(MAX(sort), 0) FROM banner").Scan(&sort)
	for _, v := range files {
		sort++
		file, _ := v.Open()
		defer file.Close()
		content, _ := ioutil.ReadAll(file)
		// 寫檔案
		targetDir := conf.Settings.Common.UPLOADS_PATH + "/banner"
		var uuid string
		db.QueryRow("SELECT uuid_short()").Scan(&uuid)
		fileName := uuid + ".jpg"
		service.WriteFile(targetDir, fileName, content)
		_, err := db.Exec(`INSERT INTO banner (picture, sort) VALUES (?, ?)`, fileName, sort)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err,
			})
			log.Panic(err)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "success",
	})
}

func UpdateBanner(c *gin.Context) {
	edit_id := c.PostForm("edit_id")
	formFile, _ := c.FormFile("file")
	file, _ := formFile.Open()
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	// 刪除舊檔案
	var pictureName string
	db.QueryRow(`SELECT picture FROM banner WHERE id=?`, edit_id).Scan(&pictureName)
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/banner"
	os.Remove(targetDir + "/" + pictureName)
	// 更新檔案
	var uuid string
	db.QueryRow("SELECT uuid_short()").Scan(&uuid)
	fileName := uuid + ".jpg"
	service.WriteFile(targetDir, fileName, content)
	_, err := db.Exec(`UPDATE banner SET picture = ? WHERE id = ?`, fileName, edit_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err,
		})
		log.Panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "success",
	})
}

func DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	// 刪除舊檔案
	var bannerPicture string
	db.QueryRow(`SELECT picture FROM banner WHERE id=?`, id).Scan(&bannerPicture)
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/banner"
	os.Remove(targetDir + "/" + bannerPicture)
	sql := fmt.Sprintf("DELETE FROM banner WHERE id=%s", id)
	_, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err,
		})
		log.Panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "success",
	})
}
