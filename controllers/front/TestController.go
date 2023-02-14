package front

import (
	"database/sql"
	"fmt"
	"goWeb/database"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func init() {
	db = database.DbConnect()
}

func Test(c *gin.Context) {
	c.HTML(http.StatusOK, "front-test", gin.H{})
}

func Test2(c *gin.Context) {
	form, _ := c.MultipartForm()
	pid := c.PostForm("pid")
	files := form.File["files[]"]
	fmt.Printf("files: %v\n", files)
	for _, v := range files {
		file, _ := v.Open()
		defer file.Close()
		content, _ := ioutil.ReadAll(file)
		sql := fmt.Sprintf("INSERT INTO products_picture (pid, picture) VALUES (%s, ?)", pid)
		_, err := db.Exec(sql, content)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}
}

func Test3(c *gin.Context) {
	sql := "SELECT id, picture FROM products_picture WHERE id=(SELECT MAX(id) FROM products_picture)"
	var id int
	var picture []byte
	db.QueryRow(sql).Scan(&id, &picture)
	fmt.Printf("id: %v\n", id)
	fmt.Printf("picture: %v\n", picture)
	c.String(200, string(picture))
}
