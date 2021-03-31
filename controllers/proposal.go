package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"tfg/aws"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

type getSignedUrlResponse struct {
	Url string `json:"url"`
}
type postSubmissionRequest struct {
	FileName   string `json:"file_name"`
	UserID     uint   `json:"user_id"`
	ProposalID uint   `json:"proposal_id"`
}

// GetProposal gets proposal from DB
func GetProposal() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := models.Proposal{}
		responseUser := models.ProposalUser{}
		_, exists := c.Get("id") // Check if auth. We may want to remove it
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
		r := database.GlobalDB.Where("id = ?", proposalIDInt).Find(&response)
		if r.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Proposal not found",
			})
			c.Abort()
			return
		}
		t := database.GlobalDB.Model(response).Association("User").Find(&responseUser)
		if t != nil {
			fmt.Println(t)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		response.User = responseUser
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
		result := database.GlobalDB.Offset(int(from)).Limit(int(to)).Preload("User").Find(&proposals)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(200, proposals)
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
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Name, description, rate and limit are required",
			})
			c.Abort()
			fmt.Println(err)
			return
		}
		intID, err := strconv.ParseUint(userID.(string), 10, 32)
		if intID != uint64(proposal.UserID) {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		user := &models.User{}
		user.ID = uint(intID)
		r := database.GlobalDB.Where(&user).Find(&proposal.User)
		err = proposal.CreateProposalRecord()
		if err != nil || r.Error != nil {
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
func GetProposalSignedUpload(session aws.StorageSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		file := c.Query("file_name")
		if file == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "file_name param required",
			})
			c.Abort()
			return
		}
		fileName, err := session.GenerateFileName(file, "pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		url, err := session.GetPutSignedUrl(fileName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, &getSignedUrlResponse{Url: url})

	}
}

func PostProposalSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		id := userID.(string)
		request := &postSubmissionRequest{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "file_name param required",
			})
			c.Abort()
			return
		}

		submission := &models.Submission{UserID: request.UserID,
			ProposalID: request.ProposalID,
			FileName:   request.FileName}
		if id != submission.UserIDString() {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		if err = submission.CreateSubmissionRecord(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, submission)
	}
}
