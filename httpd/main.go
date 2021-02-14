package main

import (
	"database/sql"
	"fmt"
	"log"
	"tfg-backend/docker"

	_ "github.com/go-sql-driver/mysql"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", credentials.Username, credentials.Password, credentials.Hostname, credentials.DbName)
}

func main() {
	docker.Compose()
	db, err := sql.Open("mysql", dsn(""))
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var version string

	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(version)
}
