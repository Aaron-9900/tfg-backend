package db_access

import (
	"fmt"
	"tfg/cmd/db"
	"tfg/cmd/db/model"
	"tfg/cmd/db/utils"
)

// GetUsers selects all users
func GetUsers() []model.PublicUser {
	users := []model.PublicUser{}
	db.GlobalDB.Model(&model.User{}).Select("name", "email", "id").Find(&users)
	fmt.Println(users)
	return users
}

func GetUser(name string) model.PublicUser {
	user := model.PublicUser{}
	db.GlobalDB.Model(&model.User{Name: name}).Select("name", "email").Find(&user)
	return user
}

func PostUser(user model.User) (model.User, error) {
	user.Password = utils.HashAndSalt(user.Password)
	result := db.GlobalDB.Create(&user)
	return user, result.Error
}
