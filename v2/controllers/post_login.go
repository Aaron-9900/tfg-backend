package controllers

import (
	"fmt"
	"log"
	"net/http"
	"tfg/v2/auth"
	"tfg/v2/database"
	"tfg/v2/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TokenResponse token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login logs users in
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload LoginPayload
		var user models.User

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Fields 'email' and 'password' required",
			})
			c.Abort()
			return
		}

		result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "invalid user credentials",
			})
			c.Abort()
			return
		}

		err = user.CheckPassword(payload.Password)
		if err != nil {
			log.Println(err)
			c.JSON(401, gin.H{
				"msg": "invalid user credentials",
			})
			c.Abort()
			return
		}
		signedToken, refreshToken, err := auth.GenerateTokens(fmt.Sprint(user.ID))
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}
		if err = database.RDB.SetUserAndToken(user.ID, refreshToken); err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}

		tokenResponse := TokenResponse{
			AccessToken:  signedToken,
			RefreshToken: refreshToken,
		}

		c.JSON(200, tokenResponse)

		return
	}

}
