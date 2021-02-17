// controllers/protected.go

package controllers

import (
	"strconv"
	"tfg/v2/database"
	"tfg/v2/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Profile returns user data
func Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		id, _ := c.Get("id") // from the authorization middleware
		intID, _ := strconv.ParseUint(id.(string), 10, 32)
		result := database.GlobalDB.Where("id = ?", intID).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"msg": "user not found",
			})
			c.Abort()
			return
		}

		if result.Error != nil {
			c.JSON(500, gin.H{
				"msg": "could not get user profile",
			})
			c.Abort()
			return
		}

		c.JSON(200, ProfileResponse{
			Name:  user.Name,
			Email: user.Email,
		})

		return
	}

}
