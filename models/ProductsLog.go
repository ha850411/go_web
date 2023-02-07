package models

import (
	"time"
)

type ProductsLog struct {
	Id         int
	Pid        int
	Amount     int
	UpdateTime time.Time
}
