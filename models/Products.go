package models

import (
	"time"
)

type Products struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Amount       int       `json:"amount"`
	AmountNotice int       `json:"amountNotice"`
	Status       int       `json:"status"`
	Type         int       `json:"type"`
	Tname        string    `json:"tname"`
	UpdateTime   time.Time `json:"updateTime"`
	FormatTime   string    `json:"formatTime"`
	PictureCnt   int       `json:"pictureCnt"`
}
