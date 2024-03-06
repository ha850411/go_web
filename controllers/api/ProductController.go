package api

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"goWeb/conf"
	"goWeb/database"
	"goWeb/models"
	"goWeb/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func init() {
	db = database.DbConnect()
}

func GetProductsList(c *gin.Context) {
	// params
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數
	keyword := c.Query("search[value]")     // 關鍵字
	where := "1"
	if keyword != "" {
		where = fmt.Sprintf("a.name LIKE '%%%s%%' OR b.name LIKE '%%%s%%'", keyword, keyword)
	}
	// 定義排序欄位
	// columes := map[string]string{
	// 	"0": "name",
	// 	"1": "amount",
	// 	"2": "amountNotice",
	// 	"3": "price",
	// 	"4": "type",
	// 	"5": "updateTime",
	// }
	// orderBy := c.DefaultQuery("order[0][column]", "0")   // 分頁筆數
	// orderType := c.DefaultQuery("order[0][dir]", "desc") // 起始筆數
	// count total
	var ids string
	sql := fmt.Sprintf(`SELECT GROUP_CONCAT( DISTINCT a.id) as ids
	FROM products as a
	LEFT JOIN products_type as b ON a.type = b.id
	WHERE %s AND a.status=1 `, where)
	db.QueryRow(sql).Scan(&ids)
	count := len(strings.Split(ids, ","))
	// 分頁
	sql = fmt.Sprintf(`SELECT a.id, a.name, a.amount, a.amountNotice, a.price, a.discount_price, a.updateTime, a.type, IFNULL(b.name, '') as tname, count(c.id) as pictureCnt, a.expiredDate
	FROM products as a
	LEFT JOIN products_type as b ON a.type = b.id
	LEFT JOIN products_picture as c ON a.id=c.pid
	WHERE a.id IN (%s) GROUP BY a.id ORDER BY a.sort asc LIMIT %s OFFSET %s`, ids, limit, offset)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := models.Products{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Amount, &rowData.AmountNotice, &rowData.Price, &rowData.DiscountPrice, &rowData.UpdateTime, &rowData.Type, &rowData.Tname, &rowData.PictureCnt, &rowData.ExpiredDate)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")

		if rowData.ExpiredDate.Unix() > 0 {
			rowData.FormatExpiredDate = rowData.ExpiredDate.Format("2006-01-02")
		}

		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func UpdateAmount(c *gin.Context) {
	id := c.PostForm("id")
	amount := c.PostForm("amount")
	sql := fmt.Sprintf("UPDATE products SET amount=%s WHERE id=%s", amount, id)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowsAffected, _ := rows.RowsAffected()
	fmt.Printf("sql: %v\n", sql)
	fmt.Printf("rowsAffected: %v\n", rowsAffected)
	if rowsAffected > 0 {
		writeProductsLog(id, amount) // 寫 log
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"updatedRows": rowsAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func AddProduct(c *gin.Context) {
	name := c.PostForm("name")
	amount := c.PostForm("amount")
	amountNotice := c.DefaultPostForm("amount_notice", "0")
	if len(amountNotice) == 0 {
		amountNotice = "0"
	}
	productType := c.PostForm("type")
	price := c.PostForm("price")
	discount_price := c.PostForm("discount_price")
	if discount_price == "" {
		discount_price = "0"
	}
	content := c.DefaultPostForm("content", "")
	expiredDate := c.DefaultPostForm("expiredDate", "")

	sql := fmt.Sprintf("INSERT INTO products (name, amount, amountNotice, type, price, discount_price, content, expiredDate) VALUES ('%s', %s, %s, %s, %s, %s, '%s', '%s')",
		name, amount, amountNotice, productType, price, discount_price, content, expiredDate)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	lastInsertID, insertErr := row.LastInsertId()
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  insertErr,
		})
		log.Panic(insertErr)
	}
	writeProductsLog(strconv.FormatInt(lastInsertID, 10), amount)
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"insertId": lastInsertID,
	})
}

func EditProduct(c *gin.Context) {
	editId := c.PostForm("id")
	name := c.PostForm("name")
	amount := c.PostForm("amount")
	amountNotice := c.DefaultPostForm("amount_notice", "0")
	if len(amountNotice) == 0 {
		amountNotice = "0"
	}
	productType := c.PostForm("type")
	price := c.PostForm("price")
	discount_price := c.DefaultPostForm("discount_price", "0")
	if len(discount_price) == 0 {
		discount_price = "0"
	}
	content := c.DefaultPostForm("content", "")

	expiredDate := conf.GetStringOrNil(c, "expiredDate")

	// 檢查數量是否有變
	var nowAmount string
	sql := fmt.Sprintf("SELECT amount FROM products WHERE id=%s", editId)
	db.QueryRow(sql).Scan(&nowAmount)

	sql = `UPDATE products
	SET name = ?, amount = ? , amountNotice = ?, type = ?, price = ?, discount_price = ?,
	content = ?, expiredDate = ? WHERE id = ?`

	fmt.Printf("sql: %v\n", sql)
	row, err := db.Exec(sql, name, amount, amountNotice, productType, price, discount_price, content, expiredDate, editId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := row.RowsAffected()
	if rowAffected > 0 {
		if amount != nowAmount {
			writeProductsLog(editId, amount)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("id: %v\n", id)
	sql := fmt.Sprintf("UPDATE products SET status=0 WHERE id=%s", id)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := row.RowsAffected()
	if rowAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}

// 寫 log
func writeProductsLog(id string, amount string) {
	today := time.Now().Format("2006-01-02")
	var check int
	sql := fmt.Sprintf("SELECT count(*) FROM products_log WHERE pid=%s AND updateDate='%s'", id, today)
	db.QueryRow(sql).Scan(&check)
	if check > 0 {
		sql = fmt.Sprintf("UPDATE products_log SET amount=%s WHERE pid=%s AND updateDate='%s'", amount, id, today)
		db.Exec(sql)
	} else {
		sql = fmt.Sprintf("INSERT INTO products_log (pid, amount, updateDate) VALUES (%s, %s, '%s')", id, amount, today)
		db.Exec(sql)
	}
}

func GetTips(c *gin.Context) {
	// 顯示存貨 < 警戒線 或 有效期限小於三個月的資料
	sql := `
	SELECT id, name, amount, amountNotice, updateTime, expiredDate
	FROM products
	WHERE status = 1
	AND amountNotice > 0 AND amount<=amountNotice
	UNION
	SELECT id, name, amount, amountNotice, updateTime, expiredDate
	FROM products
	WHERE status = 1
	AND expiredDate != '0000-00-00' 
	AND expiredDate IS NOT NULL
	AND expiredDate < NOW() + INTERVAL 3 MONTH
	`

	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := models.Products{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.Amount, &rowData.AmountNotice, &rowData.UpdateTime, &rowData.ExpiredDate)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		if rowData.ExpiredDate.Unix() > 0 {
			rowData.FormatExpiredDate = rowData.ExpiredDate.Format("2006-01-02")
		}
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

func GetProductsNameList(c *gin.Context) {
	simple := c.DefaultQuery("simple", "")
	var sql string
	if simple == "Y" {
		sql = "SELECT id, name, 0, 0 FROM products WHERE status=1 ORDER BY sort asc"
	} else {
		sql = `SELECT a.id, a.name, count(b.id) as cnt, a.price
		FROM products as a 
		LEFT JOIN products_picture as b ON a.id=b.pid 
		WHERE a.status=1 
		GROUP BY a.id ORDER BY a.sort asc`
	}

	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	data := make([]interface{}, 0)
	// return struct
	type Result struct {
		Id          int    `json:"id"`
		Name        string `json:"text"`
		PicturesCnt int    `json:"picturesCnt"`
		Price       int    `json:"price"`
	}
	for rows.Next() {
		product := Result{}
		rows.Scan(&product.Id, &product.Name, &product.PicturesCnt, &product.Price)
		data = append(data, product)
	}
	c.JSON(http.StatusOK, gin.H{
		"results": data,
		"pagination": gin.H{
			"more": true,
		},
	})
}

func GetProductsLog(c *gin.Context) {
	pid := c.Query("pid")
	month := c.Query("month")
	monthInt, _ := strconv.Atoi(month)
	monthAdd := 0 - monthInt
	// 查詢條件 預設查前 1 個月
	date := time.Now().AddDate(0, monthAdd, 0).Format("2006-01-02")
	sql := fmt.Sprintf("SELECT id, amount, updateDate FROM products_log WHERE pid=%s AND updateDate >= '%s' ORDER BY updateDate asc", pid, date)
	fmt.Printf("sql: %v\n", sql)

	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close() // close connection
	// return struct
	type Result struct {
		Id         int       `json:"id"`
		Amount     int       `json:"amount"`
		UpdateDate time.Time `json:"updateDate"`
		FormatDate string    `json:"formatDate"`
	}
	x_Value := []string{}
	y_Value := []int{}

	for rows.Next() {
		productsLog := Result{}
		rows.Scan(&productsLog.Id, &productsLog.Amount, &productsLog.UpdateDate)
		productsLog.FormatDate = productsLog.UpdateDate.Format("2006-01-02")
		x_Value = append(x_Value, productsLog.FormatDate)
		y_Value = append(y_Value, productsLog.Amount)
	}
	// 檢查是否有今天的庫存, 若沒有則塞一筆進去
	hasToday := false
	today := time.Now().Format("2006-01-02")
	for _, date := range x_Value {
		if today == date {
			hasToday = true
			break
		}
	}
	if !hasToday {
		var todayAmount int
		sql = fmt.Sprintf("SELECT amount FROM products WHERE id=%s", pid)
		db.QueryRow(sql).Scan(&todayAmount)
		x_Value = append(x_Value, today)
		y_Value = append(y_Value, todayAmount)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"pid":     pid,
			"x_Value": x_Value,
			"y_Value": y_Value,
		},
	})
}

func ExportCsv(c *gin.Context) {
	var sql string
	pid := c.DefaultQuery("pid", "0")
	month := c.DefaultQuery("month", "1")
	mode := c.Query("mode")
	productList := make(map[string]map[string]string)
	type ProductList struct {
		Id     string
		Name   string
		Amount string
	}
	// 日期區間
	startMonth, _ := strconv.Atoi(month)
	startDate := time.Now().AddDate(0, 0-startMonth, 0).Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, 0).Format("2006-01-02")
	if mode == "single" {
		// 取得商品名稱 map
		temp := ProductList{}
		sql = fmt.Sprintf("SELECT id, name, amount FROM products WHERE id=%s AND status=1", pid)
		db.QueryRow(sql).Scan(&temp.Id, &temp.Name, &temp.Amount)
		productList[temp.Id] = map[string]string{
			"name":   temp.Name,
			"amount": temp.Amount,
		}
		// 取得 log
		sql = fmt.Sprintf("SELECT pid, amount, updateDate FROM products_log WHERE pid=%s AND updateDate >'%s' AND updateDate<'%s' ORDER BY updateDate", pid, startDate, endDate)
	} else {
		// 取得商品名稱 map
		sql = "SELECT id, name, amount FROM products WHERE status=1"
		rows, err := db.Query(sql)
		if err != nil {
			log.Panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			temp := ProductList{}
			rows.Scan(&temp.Id, &temp.Name, &temp.Amount)
			productList[temp.Id] = map[string]string{
				"name":   temp.Name,
				"amount": temp.Amount,
			}
		}
		sql = fmt.Sprintf("SELECT pid, amount, updateDate FROM products_log WHERE updateDate >'%s' AND updateDate<'%s' ORDER BY pid asc, updateDate asc", startDate, endDate)
	}
	fmt.Printf("productList: %v\n", productList)
	type Result struct {
		Pid        string
		Amount     string
		UpdateDate time.Time
	}
	rows, err := db.Query(sql)
	if err != nil {
		log.Panic(err)
		c.String(http.StatusOK, err.Error())
	}
	defer rows.Close()

	dataMap := make(map[string][]map[string]string)
	for rows.Next() {
		temp := Result{}
		rows.Scan(&temp.Pid, &temp.Amount, &temp.UpdateDate)
		dataMap[temp.Pid] = append(dataMap[temp.Pid], map[string]string{
			"pid":        temp.Pid,
			"amount":     temp.Amount,
			"updateDate": temp.UpdateDate.Format("2006-01-02"),
		})
	}
	fmt.Printf("dataMap: %v\n", dataMap)

	var dataBytes = new(bytes.Buffer)
	// 設置編碼
	dataBytes.WriteString("\xEF\xBB\xBF")
	wr := csv.NewWriter(dataBytes)
	title := []string{"統計區間", startDate + "~" + endDate}
	headList := []string{"商品名稱", "商品數量", "更新日期"}
	wr.Write(title)
	wr.Write(headList)

	for thisPid, productInfo := range productList {
		if logsAry, exist := dataMap[thisPid]; exist {
			for _, log := range logsAry {
				bodyList := []string{
					productInfo["name"],
					log["amount"],
					log["updateDate"],
				}
				wr.Write(bodyList)
			}
		}
		// 寫入今日庫存
		bodyList := []string{
			productInfo["name"],
			productInfo["amount"],
			time.Now().Format("2006-01-02"),
		}
		wr.Write(bodyList)
	}
	// 清空
	wr.Flush()
	var fileName string
	if mode == "single" {
		fileName = fmt.Sprintf("商品庫存明細_%s_%s.csv", productList[pid]["name"], time.Now().Format("20060102_150405"))
	} else {
		fileName = fmt.Sprintf("商品庫存明細_%s.csv", time.Now().Format("20060102_150405"))
	}
	c.Writer.Header().Set("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	c.String(200, dataBytes.String())
}

func GetPictureLists(c *gin.Context) {
	pid := c.Query("pid")
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數

	var count int
	sql := fmt.Sprintf("SELECT count(*) FROM products_picture WHERE pid=%s", pid)
	db.QueryRow(sql).Scan(&count)

	sql = fmt.Sprintf("SELECT id, pid, picture, updateTime FROM products_picture WHERE pid = %s LIMIT %s, %s", pid, offset, limit)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close()
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := models.ProductsPicture{}
		rows.Scan(&rowData.Id, &rowData.Pid, &rowData.Picture, &rowData.UpdateTime)
		rowData.FormatTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func ShowPicture(c *gin.Context) {
	id := c.Param("id")
	sql := fmt.Sprintf("SELECT picture FROM products_picture WHERE id=%s", id)
	var picture []byte
	db.QueryRow(sql).Scan(&picture)
	c.String(200, string(picture))
}

func UploadPic(c *gin.Context) {
	form, _ := c.MultipartForm()
	pid := c.PostForm("pid")
	files := form.File["files[]"]
	for _, v := range files {
		file, _ := v.Open()
		defer file.Close()
		content, _ := ioutil.ReadAll(file)
		// 寫檔案
		targetDir := conf.Settings.Common.UPLOADS_PATH + "/" + pid
		var uuid string
		db.QueryRow("SELECT uuid_short()").Scan(&uuid)
		fileName := uuid + ".jpg"
		service.WriteFileAndCompress(targetDir, fileName, content)
		_, err := db.Exec(`INSERT INTO products_picture (pid, picture) VALUES (?, ?)`, pid, fileName)
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

func UpdatePic(c *gin.Context) {
	edit_id := c.PostForm("edit_id")
	formFile, _ := c.FormFile("file")
	file, _ := formFile.Open()
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	// 刪除舊檔案
	var productsPic models.ProductsPicture
	db.QueryRow(`SELECT id, pid, picture FROM products_picture WHERE id=?`, edit_id).Scan(&productsPic.Id, &productsPic.Pid, &productsPic.Picture)
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/" + strconv.Itoa(productsPic.Pid)
	os.Remove(targetDir + "/" + productsPic.Picture)
	// 更新檔案
	var uuid string
	db.QueryRow("SELECT uuid_short()").Scan(&uuid)
	fileName := uuid + ".jpg"
	service.WriteFileAndCompress(targetDir, fileName, content)
	_, err := db.Exec(`UPDATE products_picture SET picture = ? WHERE id = ?`, fileName, edit_id)
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

func DeletePic(c *gin.Context) {
	id := c.Param("id")
	// 刪除舊檔案
	var productsPic models.ProductsPicture
	db.QueryRow(`SELECT pid, picture FROM products_picture WHERE id=?`, id).Scan(&productsPic.Pid, &productsPic.Picture)
	targetDir := conf.Settings.Common.UPLOADS_PATH + "/" + strconv.Itoa(productsPic.Pid)
	os.Remove(targetDir + "/" + productsPic.Picture)
	sql := fmt.Sprintf("DELETE FROM products_picture WHERE id=%s", id)
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

func GetProductType(c *gin.Context) {
	limit := c.DefaultQuery("length", "10") // 分頁筆數
	offset := c.DefaultQuery("start", "0")  // 起始筆數

	var count int
	sql := "SELECT count(*) FROM products_type WHERE status=1 "
	db.QueryRow(sql).Scan(&count)

	sql = fmt.Sprintf("SELECT id, name, createTime, updateTime FROM products_type WHERE status=1 LIMIT %s, %s", offset, limit)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		log.Panic(err)
	}
	defer rows.Close()
	data := make([]interface{}, 0)
	for rows.Next() {
		rowData := models.ProductsType{}
		rows.Scan(&rowData.Id, &rowData.Name, &rowData.CreateTime, &rowData.UpdateTime)
		rowData.FormatCreateTime = rowData.CreateTime.Format("2006-01-02 15:04:05")
		rowData.FormatUpdateTime = rowData.UpdateTime.Format("2006-01-02 15:04:05")
		data = append(data, rowData)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":            200,
		"data":            data,
		"recordsTotal":    count,
		"recordsFiltered": count,
	})
}

func AddProductType(c *gin.Context) {
	name := c.PostForm("name")
	sql := fmt.Sprintf("INSERT INTO products_type (name, status,createTime) VALUES ('%s', 1, '%s')", name, time.Now().Format("2006-01-02 15:04:05"))
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	lastInsertID, insertErr := row.LastInsertId()
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  insertErr,
		})
		log.Panic(insertErr)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"insertId": lastInsertID,
	})
}

func UpdateProductType(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	sql := fmt.Sprintf("UPDATE products_type SET name='%s' WHERE id=%s", name, id)
	rows, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowsAffected, _ := rows.RowsAffected()
	fmt.Printf("sql: %v\n", sql)
	fmt.Printf("rowsAffected: %v\n", rowsAffected)
	if rowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"updatedRows": rowsAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "無資料更新",
		})
	}
}

func DeleteProductType(c *gin.Context) {
	id := c.Param("id")
	sql := fmt.Sprintf("UPDATE products_type SET status=0 WHERE id=%s", id)
	fmt.Printf("sql: %v\n", sql)
	row, err := db.Exec(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err,
		})
		log.Panic(err)
	}
	rowAffected, _ := row.RowsAffected()
	if rowAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":        http.StatusOK,
			"rowAffected": rowAffected,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "刪除失敗",
		})
	}
}

func GetProductInfo(c *gin.Context) {
	pid := c.Param("id")
	var result models.Products
	db.QueryRow(`SELECT id, name, amount, amountNotice, price, discount_price, type, IFNULL(content, ''), expiredDate
	FROM products WHERE id= ?`, pid).Scan(&result.Id, &result.Name, &result.Amount, &result.AmountNotice, &result.Price, &result.DiscountPrice, &result.Type, &result.Content, &result.ExpiredDate)

	if result.ExpiredDate.Unix() > 0 {
		result.FormatExpiredDate = result.ExpiredDate.Format("2006-01-02")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": result,
	})
}

func SortProducts(c *gin.Context) {
	sortStr := c.PostForm("sortArr")
	sortArr := strings.Split(sortStr, ",")
	fmt.Printf("sortArr: %v\n", sortArr)

	for index, id := range sortArr {
		sort := index + 1
		sql := fmt.Sprintf("UPDATE products SET sort=%v WHERE id=%v", sort, id)
		_, err := db.Exec(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err,
			})
			log.Panic(err)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
