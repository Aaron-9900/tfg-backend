package main

import (
	"tfg/cmd/db"
	"tfg/cmd/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := db.Init()
	DB, _ := db.DB()
	defer DB.Close()
	http.Run(db)
}
