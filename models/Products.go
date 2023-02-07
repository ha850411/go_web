package models

import (
	"time"
)

type Products struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Amount       int       `json:"amount"`
	AmountNotice int       `json:"amountNotice"`
	UpdateTime   time.Time `json:"updateTime"`
	FormatTime   string    `json:"formatTime"`
}
