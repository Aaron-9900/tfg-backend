// controllers/public.go

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

// LoginPayload login body
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
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

		c.JSON(200, UserResponse{Name: user.Name, Email: user.Email, ID: user.ID})
	}
}
