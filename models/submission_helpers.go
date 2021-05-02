package models

import (
	"database/sql/driver"
	"fmt"
	"tfg/database"
)

func (s *Submission) CreateSubmissionRecord() error {
	result := database.GlobalDB.Create(&s)
	if result.Error != nil {
		return result.Error
	}
	result = database.GlobalDB.Preload("User").Find(&s)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type SubmissionStatus string

func (e *SubmissionStatus) Scan(value interface{}) error {
	*e = SubmissionStatus(value.([]byte))
	return nil
}

func (e SubmissionStatus) Value() (driver.Value, error) {
	return string(e), nil
}

const (
	Pending  SubmissionStatus = "pending"
	Rejected SubmissionStatus = "rejected"
	Accepted SubmissionStatus = "accepted"
)

// IDString returns a proposal's ID as string
func (s *Submission) IDString() string {
	return fmt.Sprint(s.ID)
}

func (s *Submission) UserIDString() string {
	return fmt.Sprint(s.UserID)

}
