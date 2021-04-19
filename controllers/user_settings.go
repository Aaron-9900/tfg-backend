package controllers

import (
	"net/http"
	"strconv"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

type putUserSettings struct {
	PrivacyPolicy string `json:"privacy_policy"`
}
type privacyPolicyResponse struct {
	models.User
	PrivacyPolicy string `json:"privacy_policy"`
}

func PutUserSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}
		request := &putUserSettings{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "submission_id, proposal_id, and status are required. Status can only be: 'pending', 'complete', 'rejected'",
			})
			c.Abort()
			return
		}
		user := models.User{}
		id, err := strconv.ParseUint(userID.(string), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		user.ID = uint(id)
		tx := database.GlobalDB.Find(&user)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		user.PrivacyPolicy = request.PrivacyPolicy
		tx = database.GlobalDB.Save(&user)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		user.Password = ""
		response := privacyPolicyResponse{User: user, PrivacyPolicy: request.PrivacyPolicy}
		c.JSON(http.StatusOK, response)

	}
}

func GetUserDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("user_id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "user_id is required",
			})
			c.Abort()
			return
		}
		id, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		user := models.User{}
		user.ID = uint(id)
		tx := database.GlobalDB.Find(&user)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		user.Password = ""
		response := privacyPolicyResponse{User: user, PrivacyPolicy: user.PrivacyPolicy}
		c.JSON(http.StatusOK, response)

	}
}
