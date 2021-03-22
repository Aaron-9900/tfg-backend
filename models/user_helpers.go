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

// IDString returns a user's ID as a string
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

// GetUserProposals returns the proposals belonging to a user
func (user *User) GetUserProposals() ([]Proposal, error) {
	proposals := []Proposal{}
	result := database.GlobalDB.Where("user_id = ?", user.ID).Find(&proposals)
	if result.Error != nil {
		return []Proposal{}, result.Error
	}
	return proposals, nil
}

// IDString returns a proposal user's ID as string
func (user *ProposalUser) IDString() string {
	return fmt.Sprint(user.ID)
}

// TableName changes the name of the table so gorm knows what to look for
func (user ProposalUser) TableName() string {
	return "users"
}
