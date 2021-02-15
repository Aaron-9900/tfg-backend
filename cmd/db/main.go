package db

import (
	"fmt"
	"tfg/credentials"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", credentials.Username, credentials.Password, credentials.Hostname, credentials.Dbname)
}

// Init initializes database. Returns database instance
func Init() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
