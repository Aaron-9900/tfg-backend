package models

import (
	"fmt"
	"tfg/database"
)

// CreateProposalRecord creates a proposal record in the database
func (p *Proposal) CreateProposalRecord() error {
	result := database.GlobalDB.Create(&p)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreateProposalRecord creates a proposal record in the database
func (p *Proposal) IDString() string {
	return fmt.Sprint(p.ID)
}
