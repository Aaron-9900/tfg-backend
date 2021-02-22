package models

import (
	"fmt"
	"tfg/database"

	"golang.org/x/crypto/bcrypt"
)

// CreateUserRecord creates a user record in the database
func (user *User) CreateUserRecord() error {
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// HashPassword encrypts user password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)

	return nil
}
func (user *User) IDString() string {
	return fmt.Sprint(user.ID)
}

// CheckPassword checks user password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (user *User) GetUserProposals() ([]Proposal, error) {
	proposals := []Proposal{}
	result := database.GlobalDB.Where("user_id = ?", user.ID).Find(&proposals)
	if result.Error != nil {
		return []Proposal{}, result.Error
	}
	return proposals, nil
}
func (user *ProposalUser) IDString() string {
	return fmt.Sprint(user.ID)
}

func (u ProposalUser) TableName() string {
	return "users"
}
