package config

import (
	"encoding/json"
	"log"
	"os"

	"api-boilerplate/models/domain"
)

var Config domain.Cfg

func InitConfig() {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
}
