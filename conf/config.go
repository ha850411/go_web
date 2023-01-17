package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	ServerPort string `yaml:"ServerPort"`
	Host       string `yaml:"Host"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	Dbname     string `yaml:"Dbname"`
	Port       string `yaml:"Port"`
}

func GetConfig() Conf {
	var settings Conf
	yamlFile, err := ioutil.ReadFile("./conf/config.yaml")
	if err != nil {
		fmt.Printf("read fail: %s", err)
	}
	//yaml文件内容影射到结构体中
	err1 := yaml.Unmarshal(yamlFile, &settings)
	if err1 != nil {
		fmt.Println("error")
	}
	return settings
}
