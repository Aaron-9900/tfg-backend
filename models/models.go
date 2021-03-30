// models/models.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// GenericDbData is a copy of gorm.Model so we can name json fields
type GenericDbData struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProposalUser defines the user json returned in a proposal
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
	Rate        float32      `json:"rate"`
	Type        string       `json:"type"`
}

type ProposalType struct {
	GenericDbData
	Value string `json:"value"`
}
