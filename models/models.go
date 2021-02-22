// models/models.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type GenericDbData struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProposalUser struct {
	GenericDbData
	Name     string `json:"name"`
	Email    string `json:"-" gorm:"unique"`
	Password string `json:"-"`
}

// User defines the user in db
type User struct {
	GenericDbData
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// Proposal defines proposal in db
type Proposal struct {
	GenericDbData
	User        ProposalUser `json:"user"`
	UserID      uint         `json:"user_id"`
	Limit       int          `json:"limit"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
}
