package service

import (
	"goWeb/database"
	"goWeb/models"
	"log"

	"github.com/goccy/go-json"
)

func GetMenu() string {
	database := database.DbConnect()
	sql := "SELECT name, icon, path, active FROM menu ORDER BY sort asc"
	rows, err := database.Query(sql)
	if err != nil {
		log.Panic("get menu fail" + err.Error())
	}
	defer rows.Close()
	menuData := []map[string]string{}
	for rows.Next() {
		menu := models.Menu{}
		rows.Scan(&menu.Name, &menu.Icon, &menu.Path, &menu.Active)
		menuData = append(menuData, map[string]string{
			"name":   menu.Name,
			"icon":   menu.Icon,
			"path":   menu.Path,
			"active": menu.Active,
		})
	}
	jsonStr, _ := json.Marshal(menuData)
	return string(jsonStr)
}
