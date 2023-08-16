package config

import (
	"encoding/json"
	"log"
	"os"

	"api-boilerplate/models/domain"
)

func InitConfig() domain.Cfg {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}

	defer file.Close()

	config := domain.Cfg{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}

	return config
}
