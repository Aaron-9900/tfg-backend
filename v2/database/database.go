// database/database.go

package database

import (
	"fmt"
	"tfg/v2/credentials"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", credentials.Username, credentials.Password, credentials.Hostname, credentials.Dbname)
}

// Init initializes database. Returns database instance
func Init() (err error) {
	GlobalDB, err = gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return err
}
