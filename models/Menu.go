package models

import (
	"time"
)

type Menu struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Path       string    `json:"path"`
	Active     string    `json:"active"`
	Sort       int       `json:"sort"`
	UpdateTime time.Time `json:"updateTime"`
}
