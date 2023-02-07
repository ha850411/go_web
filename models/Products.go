package models

import (
	"time"
)

type Products struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Amount     int       `json:"amount"`
	UpdateTime time.Time `json:"updateTime"`
	FormatTime string    `json:"formatTime"`
}
