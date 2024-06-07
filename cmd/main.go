package main

import (
	"database/sql"
	"log"

	"github.com/akshtrikha/golang-ecomm/cmd/api"
	"github.com/akshtrikha/golang-ecomm/config"
	"github.com/akshtrikha/golang-ecomm/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
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

	// create an instance of our server
	// this delegates the actual server code and its custom implementation to other parts of the project
	// help readability and is a better approact
	server := api.NewAPIServer(":9000", db)

	// run the server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully Connected!")
}