package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connect to database server.
// https://github.com/go-sql-driver/mysql#timetime-support
// https://stackoverflow.com/a/52895312
func sqlConnect() {
	var err error
	if db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true&timeout=5s", config.Database.User, config.Database.Password, config.Database.Ip, config.Database.Port, config.Database.Name)); err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
}
