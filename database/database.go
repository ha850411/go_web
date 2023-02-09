package database

import (
	"database/sql"
	"fmt"
	"goWeb/conf"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DbInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

func DbConnect() *sql.DB {
	settings := conf.GetConfig()
	dbInfo := DbInfo{
		Host:     settings.Database.Host,
		Port:     settings.Database.Port,
		Username: settings.Database.Username,
		Password: settings.Database.Password,
		Dbname:   settings.Database.Dbname,
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbInfo.Username, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.Dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func QueryToMap(db *sql.DB, sSql string) ([]map[string]interface{}, error) {

	rows, err := db.Query(sSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values, valuesPoint := make([]interface{}, count), make([]interface{}, count)

	mData := make([]map[string]interface{}, 0)
	for rows.Next() {

		for i := 0; i < count; i++ {
			valuesPoint[i] = &values[i]
		}

		rows.Scan(valuesPoint...)

		entry := make(map[string]interface{})

		for i, val := range values {
			key := columns[i]
			var v interface{}
			b, ok := val.([]byte)
			if ok {
				// 字符切片转为字符串
				v = string(b)
			} else {
				v = val
			}
			entry[key] = v
		}
		mData = append(mData, entry)
	}
	return mData, nil
}
