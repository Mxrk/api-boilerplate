package main

import (
	"api-boilerplate/config"
	"api-boilerplate/database"
	"api-boilerplate/server"
)

func main() {
	config.InitConfig()
	database.ConnectDatabase()
	server.InitServer()
}
