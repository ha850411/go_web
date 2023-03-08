package models

import (
	"time"
)

type Banner struct {
	Id         int       `json:"id"`
	Picture    string    `json:"picture"`
	Sort       int       `json:"sort"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
