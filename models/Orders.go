package models

import (
	"time"
)

type Orders struct {
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Contact    string        `json:"contact"`
	Address    string        `json:"address"`
	Remark     string        `json:"remark"`
	UpdateTime time.Time     `json:"updateTime"`
	FormatTime string        `json:"formatTime"`
	Detail     []interface{} `json:"detail"`
	Total      int           `json:"total"`
}
