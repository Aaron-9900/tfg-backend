package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tfg/auth"
	"tfg/database"

	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh tokens
func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("id")
		if !exists {
			c.JSON(401, gin.H{
				"msg": "Incorrect token",
			})
			c.Abort()
			return
		}
		parsedID, err := strconv.ParseUint(id.(string), 10, 32)
		if err != nil {
			c.JSON(401, gin.H{
				"msg": "Incorrect token",
			})
			c.Abort()
			return
		}
		token, exists := c.Get("token")
		if !exists {
			c.JSON(401, gin.H{
				"msg": "Incorrect token",
			})
			c.Abort()
			return
		}
		val, err := database.RDB.DeleteToken(uint(parsedID), token.(string))
		if err != nil || val == int64(0) {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Invalid token",
			})
			c.Abort()
			return
		}
		signedToken, refreshToken, err := auth.GenerateTokens(fmt.Sprint(parsedID))
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}
		if err = database.RDB.SetUserAndToken(uint(parsedID), refreshToken); err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}
		tokenResponse := refreshResponse{
			AccessToken:  signedToken,
			RefreshToken: refreshToken,
		}
		c.JSON(200, tokenResponse)
	}
}
