package controllers

import (
	"encoding/json"
	"fmt"
	"goWeb/database"
	"goWeb/service/linebot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func LineBotNotify(c *gin.Context) {
	code := c.Query("code")

	// 呼叫 linebot api 取得 access_token
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", linebot.REDIRECT_URI)
	form.Add("client_id", linebot.CLIENT_ID)
	form.Add("client_secret", linebot.CLIENT_SECRET)

	req, _ := http.NewRequest("POST", linebot.OAUTH_TOKEN_URL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, _ := client.Do(req)
	response, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse map[string]interface{}
	fmt.Printf("response: %v\n", string(response))
	json.Unmarshal(response, &jsonResponse)
	fmt.Printf("jsonResponse: %v\n", jsonResponse)
	fmt.Printf("jsonResponse[\"status\"]: %T\n", jsonResponse["status"])
	fmt.Printf("jsonResponse[\"status\"]: %v\n", jsonResponse["status"] == 200)
	if jsonResponse["status"] == float64(200) {
		accessToken := jsonResponse["access_token"]
		username, _ := c.Get("username")
		db := database.DbConnect()
		sql := fmt.Sprintf(`UPDATE users SET linebot_token='%s' WHERE username='%s'`, accessToken, username)
		fmt.Printf("sql: %v\n", sql)
		rows, err := db.Exec(sql)
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `<script>alert('`+err.Error()+`');location='/admin/setting'</script>`)
		}
		if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 0 {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `<script>alert('綁定成功');location='/admin/setting'</script>`)
		} else {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `<script>alert('綁定失敗');location='/admin/setting'</script>`)
		}
	} else {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<script>alert('取得 token 失敗請重新操作');location='/admin/setting'</script>`)
	}
}
