package models

import (
	"time"
)

type OrdersDetail struct {
	Id         int       `json:"id"`
	OrderId    int       `json:"order_id"`
	Pid        int       `json:"pid"`
	Pname      string    `json:"pname"`
	Amount     int       `json:"amount"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
