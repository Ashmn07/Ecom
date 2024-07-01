package main

import (
	"database/sql"
	"log"

	"github.com/Ashmn07/Ecom/cmd/api"
	"github.com/Ashmn07/Ecom/config"
	"github.com/Ashmn07/Ecom/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMyQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully Connected")
}
