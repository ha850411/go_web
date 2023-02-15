package models

import (
	"time"
)

type ProductsType struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Status           int       `json:"status"`
	CreateTime       time.Time `json:"createTime"`
	FormatCreateTime string    `json:"formatCreateTime"`
	UpdateTime       time.Time `json:"updateTime"`
	FormatUpdateTime string    `json:"formatUpdateTime"`
}
