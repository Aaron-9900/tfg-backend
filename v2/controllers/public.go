// controllers/public.go

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"tfg/v2/auth"
	"tfg/v2/credentials"
	"tfg/v2/database"
	"tfg/v2/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginPayload login body
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse token response
type LoginResponse struct {
	Token string `json:"token"`
}

// Signup creates a user in db
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		err := c.ShouldBindJSON(&user)
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "invalid json",
			})
			c.Abort()

			return
		}

		err = user.HashPassword(user.Password)
		if err != nil {
			log.Println(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error hashing password",
			})
			c.Abort()

			return
		}

		err = user.CreateUserRecord()
		fmt.Println(err)
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusConflict, gin.H{
				"msg": "User already exists",
			})
			c.Abort()

			return
		}

		c.JSON(200, user)
	}
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

		jwtWrapper := auth.JwtWrapper{
			SecretKey:       credentials.JwtKey,
			Issuer:          "AuthService",
			ExpirationHours: 72,
		}

		signedToken, err := jwtWrapper.GenerateToken(user.Email)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}

		tokenResponse := LoginResponse{
			Token: signedToken,
		}

		c.JSON(200, tokenResponse)

		return
	}

}
