package models

import (
	"time"
)

type ProductsPicture struct {
	Id         int       `json:"id"`
	Pid        int       `json:"pid"`
	Picture    string    `json:"picture"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
