package main

import (
	"tfg/cmd/db"
	"tfg/cmd/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := db.Init()
	if err != nil {
		panic("DB ERROR")
	}
	DB, err := db.GlobalDB.DB()
	if err != nil {
		panic("DB ERROR")
	}
	defer DB.Close()
	http.Run()
}
