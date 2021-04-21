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
type LowInfoUser struct {
	GenericDbData
	Name          string `json:"name"`
	Email         string `json:"-" gorm:"unique"`
	Password      string `json:"-"`
	PrivacyPolicy string `json:"privacy_policy,omitempty"`
}

// User defines the user in db
type User struct {
	GenericDbData
	Name          string `json:"name"`
	Email         string `json:"email,omitempty" gorm:"unique"`
	Password      string `json:"password,omitempty"`
	PrivacyPolicy string `json:"-"`
}

// Proposal defines proposal in db
type Proposal struct {
	GenericDbData
	User        LowInfoUser  `json:"user,omitempty"`
	UserID      uint         `json:"-"`
	Submissions []Submission `json:"submissions,omitempty"`
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

type Submission struct {
	GenericDbData
	UserID      uint             `json:"-"`
	User        LowInfoUser      `json:"user"`
	ProposalID  uint             `json:"-"`
	Proposal    Proposal         `json:"proposal"`
	FileName    string           `json:"file_name"`
	Status      SubmissionStatus `json:"status"`
	ContentType string           `json:"content_type"`
}

type PrivacyTemplates struct {
	GenericDbData
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
}
