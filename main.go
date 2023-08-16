package main

import (
	"log"

	"api-boilerplate/config"
	"api-boilerplate/database"
	"api-boilerplate/models/domain"
	"api-boilerplate/server"
)

type Main struct {
	Config domain.Cfg

	DB *database.DB

	HTTPServer *server.Server
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	cfg := config.InitConfig()
	return &Main{
		Config:     cfg,
		DB:         database.NewDB(cfg),
		HTTPServer: server.InitServer(),
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	m := NewMain()
	m.Run()
}

func (m *Main) Run() {
	m.HTTPServer.UserService = database.NewUserService(m.DB)

	log.Println("Server started on port " + m.Config.Server.Port)
	err := m.HTTPServer.StartServer(":" + m.Config.Server.Port)
	if err != nil {
		log.Fatal(err)
	}
}
