package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"tfg/aws"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type getSignedUrlResponse struct {
	Url      string `json:"url"`
	FileName string `json:"file_name"`
}
type getGetSignedUrlResponse struct {
	Url string `json:"url"`
}
type postSubmissionRequest struct {
	FileName    string `json:"file_name"`
	UserID      uint   `json:"user_id"`
	ProposalID  uint   `json:"proposal_id"`
	ContentType string `json:"content_type"`
}
type postSubmissionResponse struct {
	models.GenericDbData
	UserID      uint                    `json:"-"`
	User        models.LowInfoUser      `json:"user"`
	ProposalID  uint                    `json:"-"`
	Proposal    models.Proposal         `json:"-"`
	FileName    string                  `json:"file_name"`
	Status      models.SubmissionStatus `json:"status"`
	ContentType string                  `json:"content_type"`
}
type setSubmissionStatusRequest struct {
	SubmissionID uint                    `json:"submission_id"`
	Status       models.SubmissionStatus `json:"status"`
	ProposalID   uint                    `json:"proposal_id"`
}
type getUserSubmissionsRequest struct {
	ProposalID uint `json:"proposal_id"`
}
type putUserBalance struct {
	Balance float64 `json:"balance"`
}
type userBalanceResponse struct {
	models.User
	Balance float64 `json:"balance"`
}

// TODO: Add test

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
		fileName, err := session.GenerateFileName(file)
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

// TODO: Add test

func GetProposalSignedDownload(session aws.StorageSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		file := c.Query("file_name")
		submissionID := c.Query("submission_id")
		userID, exists := c.Get("id") // from the authorization middleware
		if file == "" || submissionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "file_name and submission_id params required",
			})
			c.Abort()
			return
		}
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		intID, err := strconv.ParseUint(userID.(string), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		id := uint(intID)
		submission := &models.Submission{}
		submissionIDInt, err := strconv.ParseInt(submissionID, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		submission.ID = uint(submissionIDInt)
		if tx := database.GlobalDB.Preload("Proposal").Preload("Proposal.User").Find(&submission); tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		if id != submission.UserID || id != submission.Proposal.UserID {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		url, err := session.GetGetSignedUrl(file)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, &getGetSignedUrlResponse{Url: url})

	}
}

// TODO: Add test

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
			ProposalID:  request.ProposalID,
			FileName:    request.FileName,
			Status:      models.Pending,
			ContentType: request.ContentType,
		}
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

// TODO: Add test
func GetUserSubmissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		request := c.Query("proposal_id")
		if request == "" {
			request = "0"
		}
		requestInt, err := strconv.ParseUint(request, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		userID, err = strconv.ParseUint(userID.(string), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		submissions := []models.Submission{}
		if requestInt != 0 {
			if tx := database.GlobalDB.Where("user_id = ? AND proposal_id = ?", uint(userID.(uint64)), requestInt).
				Preload("User").Preload("Proposal").Preload("Proposal.User").Order("created_at").Find(&submissions); tx.Error != nil {
				fmt.Println(tx)
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Server error",
				})
				c.Abort()
				return
			}
		} else {
			if tx := database.GlobalDB.Where("user_id = ?", uint(userID.(uint64))).Preload("User").Preload("Proposal").
				Preload("Proposal.User").Find(&submissions); tx.Error != nil {
				fmt.Println(tx)

				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Server error",
				})
				c.Abort()
				return
			}
		}
		c.JSON(http.StatusOK, submissions)

	}
}

func PostSubmissionStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		request := &setSubmissionStatusRequest{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "submission_id, proposal_id, and status are required. Status can only be: 'pending', 'complete', 'rejected'",
			})
			c.Abort()
			return
		}
		submission := &models.Submission{}
		submission.ID = request.SubmissionID
		database.GlobalDB.Preload("Proposal").Preload("Proposal.User").Find(&submission)
		if userID.(string) != submission.Proposal.User.IDString() {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		if request.Status == "accepted" {
			userSubmitting := &models.User{}
			userIDint, err := strconv.ParseUint(userID.(string), 10, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Server error",
				})
				c.Abort()
				return
			}
			userSubmitting.ID = uint(userIDint)
			tx := database.GlobalDB.Find(&userSubmitting)
			if tx.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Server error",
				})
				c.Abort()
				return
			}
			if userSubmitting.Balance-float64(submission.Proposal.Rate) < 0 {
				c.JSON(http.StatusPreconditionFailed, gin.H{
					"msg": "Not enough balance",
				})
				c.Abort()
				return
			}
			err = database.GlobalDB.Transaction(func(tx *gorm.DB) error {
				targetUser := submission.User
				targetUser.Balance = targetUser.Balance + float64(submission.Proposal.Rate)
				if r := tx.Save(&targetUser); r.Error != nil {
					return r.Error
				}
				userSubmitting.Balance = userSubmitting.Balance - float64(submission.Proposal.Rate)
				if r := tx.Save(&userSubmitting); r.Error != nil {
					return r.Error
				}
				return nil
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Server error",
				})
				c.Abort()
				return
			}

		}
		submission.Status = request.Status
		if tx := database.GlobalDB.Save(&submission); tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, submission)

	}
}

func PutUserBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		request := &putUserBalance{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "balance is required",
			})
			c.Abort()
			return
		}

		user := &models.User{}
		id, err := strconv.ParseUint(userID.(string), 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "server error",
			})
			c.Abort()
			return
		}
		user.ID = uint(id)
		if err = database.GlobalDB.Find(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "server error",
			})
			c.Abort()
			return
		}
		user.Balance = user.Balance + request.Balance
		if err = database.GlobalDB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "server error",
			})
			c.Abort()
			return
		}
		response := &userBalanceResponse{User: *user, Balance: user.Balance}
		c.JSON(http.StatusOK, response)

	}
}
