package linebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goWeb/conf"
	"net/http"
)

var (
	REDIRECT_URI    = ""
	OAUTH_TOKEN_URL = "https://notify-bot.line.me/oauth/token"
	NOTICE_URL      = "https://api.line.me/v2/bot/message/broadcast"
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
	channelAccessToken := conf.Settings.Common.LINE_MESSAGING_API_ACCESS_TOKEN
	// 創建請求體
	broadcastPayload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	// 編碼為 JSON
	jsonData, err := json.Marshal(broadcastPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// 創建 HTTP POST 請求
	req, err := http.NewRequest("POST", NOTICE_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 設置請求標頭
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 檢查響應狀態
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
	} else {
		fmt.Printf("Failed to send message, status code: %d\n", resp.StatusCode)
	}
}
