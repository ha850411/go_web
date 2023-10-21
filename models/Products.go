package models

import (
	"time"
)

type Products struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Amount            int       `json:"amount"`
	AmountNotice      int       `json:"amountNotice"`
	Price             int       `json:"price"`
	DiscountPrice     int       `json:"discount_price"`
	Status            int       `json:"status"`
	Type              int       `json:"type"`
	Tname             string    `json:"tname"`
	Content           string    `json:"content"`
	ExpiredDate       time.Time `json:"expiredDate"`
	UpdateTime        time.Time `json:"updateTime"`
	FormatExpiredDate string    `json:"formatExpiredDate"`
	FormatTime        string    `json:"formatTime"`
	PictureCnt        int       `json:"pictureCnt"`
}
