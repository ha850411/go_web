package models

import (
	"time"
)

type Orders struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Contact    string    `json:"contact"`
	Pid        string    `json:"pid"`
	Pname      string    `json:"pname"`
	Amount     string    `json:"amount"`
	Remark     string    `json:"remark"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
