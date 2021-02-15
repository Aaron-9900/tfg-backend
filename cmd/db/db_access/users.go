package db_access

import (
	"fmt"
	"log"
	"tfg/cmd/db/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUsers selects all users
func GetUsers(db *gorm.DB) []model.PublicUser {
	users := []model.PublicUser{}
	db.Model(&model.User{}).Select("name", "email").Find(&users)
	fmt.Println(users)
	return users
}

func GetUser(db *gorm.DB, name string) model.PublicUser {
	user := model.PublicUser{}
	db.Model(&model.User{Name: name}).Select("name", "email").Find(&user)
	return user
}

func PostUser(user model.User, db *gorm.DB) (model.User, error) {
	user.Password = hashAndSalt(user.Password)
	result := db.Create(&user)
	return user, result.Error
}

func hashAndSalt(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
