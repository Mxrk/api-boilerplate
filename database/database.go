package database

import (
	"fmt"
	"log"

	"api-boilerplate/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var db *sqlx.DB

// ConnectDatabase initials the database
func ConnectDatabase() {
	var err error
	cfg := config.Config.Database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable timezone=CET",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "database/migrations",
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Applied %d migrations!\n", n)
	log.Println("Connected to database")
}
