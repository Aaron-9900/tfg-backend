package controllers

import (
	"net/http"
	"strconv"
	"tfg/database"

	"github.com/gin-gonic/gin"
)

// Login logs users in
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, exists := c.Get("id") // from the authorization middleware
		intID, _ := strconv.ParseUint(id.(string), 10, 32)

		if !exists {
			c.JSON(401, gin.H{
				"msg": "Incorrect token",
			})
			c.Abort()
			return
		}
		token, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Incorrect token",
			})
			c.Abort()
			return
		}
		deleted, err := database.RDB.DeleteToken(uint(intID), token.(string))
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "error deleting token",
			})
			c.Abort()
			return
		}
		if deleted == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Token not found",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"msg": "Success in logout",
		})

		return
	}

}
