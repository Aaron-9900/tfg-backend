package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

type getProposalType struct {
	models.Proposal
	SubmissionCount   int  `json:"submission_count" gorm:"-"`
	HasUserSubmission bool `json:"has_user_submission"`
}
type getProposalsResponse struct {
	Proposals []models.Proposal `json:"proposals"`
	Count     int64             `json:"count"`
}

// GetProposal gets proposal from DB
func GetProposal() gin.HandlerFunc {
	return func(c *gin.Context) {
		tempResponse := models.Proposal{}
		userID, exists := c.Get("id") // Check if auth. We may want to remove it
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}

		proposalID := c.Query("proposal_id")
		if proposalID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "proposal_id required",
			})
			c.Abort()
			return
		}
		proposalIDInt, err := strconv.ParseInt(proposalID, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		tempResponse.ID = uint(proposalIDInt)
		r := database.GlobalDB.Preload("Submissions").Preload("User").Preload("Submissions.User").Find(&tempResponse)
		if r.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Proposal not found",
			})
			c.Abort()
			return
		}
		response := &getProposalType{Proposal: tempResponse}
		for i := range tempResponse.Submissions {
			if tempResponse.Submissions[i].User.IDString() == userID {
				response.HasUserSubmission = true
			}
		}
		if userID != tempResponse.User.IDString() {
			database.GlobalDB.Preload("User").Where("user_id = ? AND proposal_id = ?", userID, tempResponse.IDString()).Find(&response.Submissions)
		}
		response.SubmissionCount = len(response.Submissions)
		c.JSON(200, response)
		return
	}

}

// GetProposals returns list of proposals
func GetProposals() gin.HandlerFunc {
	return func(c *gin.Context) {
		proposals := []models.Proposal{}

		fromStr := c.Query("from")
		toStr := c.Query("to")

		if fromStr == "" || toStr == "" {
			fromStr = "-1"
			toStr = "-1"
		}

		from, err := strconv.ParseInt(fromStr, 10, 64)
		to, err := strconv.ParseInt(toStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Please, provide valid 'from' and 'to' fields",
			})
			c.Abort()
			return
		}
		result := database.GlobalDB.Order("created_at desc").Offset(int(from)).Limit(int(to)).Preload("User").Find(&proposals)
		for i := range proposals {
			proposals[i].User.PrivacyPolicy = ""
		}
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		var count int64 = 0
		database.GlobalDB.Model(&models.Proposal{}).Count(&count)
		response := &getProposalsResponse{Proposals: proposals, Count: count}
		c.JSON(200, response)
		return

	}
}

// PostProposal sets proposal in DB
func PostProposal() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		proposal := models.Proposal{}
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}

		err := c.ShouldBindJSON(&proposal)
		if err != nil || proposal.Name == "" || proposal.Description == "" || proposal.Limit == 0 || proposal.Rate == 0 {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Name, description, rate and limit are required",
			})
			c.Abort()
			fmt.Println(err)
			return
		}
		intID, err := strconv.ParseUint(userID.(string), 10, 32)

		user := &models.User{}
		user.ID = uint(intID)
		t := database.GlobalDB.Find(&user)
		if t.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		if user.PrivacyPolicy == "" {
			c.JSON(http.StatusPreconditionFailed, gin.H{
				"msg": "User needs a privacy policy",
			})
			c.Abort()
			return
		}
		proposal.UserID = uint(intID)
		err = proposal.CreateProposalRecord()
		t = database.GlobalDB.Where("id = ?", user.ID).Find(&proposal.User)
		if err != nil || t.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		c.JSON(200, proposal)

		return
	}

}
func GetProposalTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		proposalTypes := []models.ProposalType{}
		result := database.GlobalDB.Find(&proposalTypes)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, proposalTypes)

	}
}
