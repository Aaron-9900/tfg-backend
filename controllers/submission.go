package controllers

import (
	"net/http"
	"tfg/aws"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

type getSignedUrlResponse struct {
	Url      string `json:"url"`
	FileName string `json:"file_name"`
}
type postSubmissionRequest struct {
	FileName   string `json:"file_name"`
	UserID     uint   `json:"user_id"`
	ProposalID uint   `json:"proposal_id"`
}
type postSubmissionResponse struct {
	models.GenericDbData
	UserID     uint               `json:"-"`
	User       models.LowInfoUser `json:"user"`
	ProposalID uint               `json:"-"`
	Proposal   models.Proposal    `json:"-"`
	FileName   string             `json:"file_name"`
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
		c.JSON(http.StatusOK, &getSignedUrlResponse{Url: url, FileName: fileName})

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
		response := postSubmissionResponse(*submission)
		c.JSON(http.StatusOK, response)
	}
}
