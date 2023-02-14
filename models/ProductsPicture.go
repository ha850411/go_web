package models

import (
	"time"
)

type ProductsPicture struct {
	Id         int       `json:"id"`
	Pid        int       `json:"pid"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
