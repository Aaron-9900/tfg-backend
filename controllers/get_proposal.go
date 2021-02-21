package controllers

import (
	"net/http"
	"strconv"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

// GetProposal gets proposal from DB
func GetProposal() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := models.Proposal{}
		_, exists := c.Get("id") // Check if auth. We may want to remove it
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}

		proposalID := c.Query("proposal_id")

		if proposalID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "proposal_id required",
			})
			c.Abort()
			return
		}
		proposalIDInt, err := strconv.ParseInt(proposalID, 10, 64)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		r := database.GlobalDB.Where("id = ?", proposalIDInt).Find(&response)
		if r.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Proposal not found",
			})
			c.Abort()
			return
		}
		if r.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		c.JSON(200, response)

		return
	}

}
