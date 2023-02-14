package controllers

import (
	"fmt"
	"goWeb/database"
	"goWeb/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product") // 取得 menu
	keyword := c.DefaultQuery("keyword", "")
	output["keyword"] = keyword
	c.HTML(http.StatusOK, "product", output)
}

func ProductPicManage(c *gin.Context) {
	output := service.GetCommonOutput(c, "product") // 取得 menu
	pid := c.Param("pid")
	output["pid"] = pid
	output["pname"] = getProductName(pid)
	fmt.Printf("pid: %v\n", pid)
	c.HTML(http.StatusOK, "productPic", output)
}

func getProductName(id string) (pname string) {
	db := database.DbConnect()
	db.QueryRow(`SELECT name FROM products WHERE id=?`, id).Scan(&pname)
	return
}
