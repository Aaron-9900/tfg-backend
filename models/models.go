// models/models.go

package models

import (
	"gorm.io/gorm"
)

// User defines the user in db
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// Proposal defines proposal in db
type Proposal struct {
	gorm.Model
	User        User   `json:"-"`
	UserID      uint   `json:"user_id"`
	Limit       int    `json:"limit"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
