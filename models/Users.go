package models

import "database/sql"

type Users struct {
	Id         int
	Username   string
	Password   string
	UpdateTime sql.NullTime
}
