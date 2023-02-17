package linebot

import (
	"fmt"
	"goWeb/conf"
	"goWeb/database"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	REDIRECT_URI    = ""
	OAUTH_TOKEN_URL = "https://notify-bot.line.me/oauth/token"
	NOTICE_URL      = "https://notify-api.line.me/api/notify"
	CLIENT_ID       = "Hg6jcH8P2quWxvlBEpC4n8"
	CLIENT_SECRET   = "YE5r6rhpBFUGjqyeK7z9TCUKsSkxK2DrF0EpZAjntUM"
)

func init() {
	if conf.Settings.Common.ENV == "local" {
		REDIRECT_URI = "http://localhost:8087/linebot/notify"
	} else {
		REDIRECT_URI = "http://107.191.52.212/linebot/notify"
	}
}

func Request(message string) {
	form := url.Values{}
	form.Add("message", message)

	sql := "SELECT linebot_token FROM users WHERE linebot_token IS NOT NULL"
	db := database.DbConnect()
	rows, err := db.Query(sql)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var token string
		rows.Scan(&token)
		bearerToken := fmt.Sprintf("Bearer %s", token)
		req, _ := http.NewRequest("POST", NOTICE_URL, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", bearerToken)

		client := &http.Client{}
		resp, _ := client.Do(req)
		response, _ := ioutil.ReadAll(resp.Body)
		// var jsonResponse map[string]interface{}
		fmt.Printf("response: %v\n", string(response))
	}
}
