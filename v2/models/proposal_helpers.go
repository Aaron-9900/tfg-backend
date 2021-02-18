package models

import "tfg/v2/database"

// CreateProposalRecord creates a proposal record in the database
func (p *Proposal) CreateProposalRecord() error {
	result := database.GlobalDB.Create(&p)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
