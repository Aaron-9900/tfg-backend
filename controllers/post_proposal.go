package controllers

import (
	"net/http"
	"strconv"
	"tfg/models"

	"github.com/gin-gonic/gin"
)

// Proposal sets proposal in DB
func PostProposal() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id") // from the authorization middleware
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Forbidden",
			})
			c.Abort()
			return
		}

		name := c.Query("name")
		description := c.Query("description")
		limit := c.Query("limit")

		if name == "" || description == "" || limit == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Missing parameters. name, description and limit are required",
			})
			c.Abort()
			return
		}

		intID, err := strconv.ParseUint(userID.(string), 10, 32)
		limitInt, err := strconv.ParseInt(limit, 10, 64)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		proposal := models.Proposal{UserID: uint(intID), Name: name, Description: description, Limit: int(limitInt)}

		err = proposal.CreateProposalRecord()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}

		c.JSON(200, proposal)

		return
	}

}
