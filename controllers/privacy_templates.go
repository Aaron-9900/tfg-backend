package controllers

import (
	"net/http"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

func GetPrivacyTemplates() gin.HandlerFunc {
	return func(c *gin.Context) {
		templates := []models.PrivacyTemplates{}
		response := database.GlobalDB.Find(&templates)

		if response.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(200, templates)
		return

	}
}
