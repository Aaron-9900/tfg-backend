package repo

import (
	"tfg/models"

	"gorm.io/gorm"
)

type repository interface {
	GetUserById(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GenerateTokens(id string) (string, string, error)
	SetTokens(id string, token string) error
}

type newRepository struct {
	database gorm.DB
}

func (r *newRepository) GetUserById(id string) (models.User, error) {
	user := models.User{}
	err := r.database.Where("id = ?", id).First(&user)
	return user, err.Error
}

func InitRepo() {
	// return &newRepository{database: database.Init()}
}
