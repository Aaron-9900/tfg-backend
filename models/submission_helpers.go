package models

import (
	"fmt"
	"tfg/database"
)

func (s *Submission) CreateSubmissionRecord() error {
	result := database.GlobalDB.Create(&s)
	if result.Error != nil {
		return result.Error
	}
	result = database.GlobalDB.Preload("Proposal").Preload("User").Preload("Proposal.User").Find(&s)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IDString returns a proposal's ID as string
func (s *Submission) IDString() string {
	return fmt.Sprint(s.ID)
}

func (s *Submission) UserIDString() string {
	return fmt.Sprint(s.UserID)

}
