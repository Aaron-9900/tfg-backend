package controllers

import (
	"fmt"
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

// GetProposals returns list of proposals
func GetProposals() gin.HandlerFunc {
	return func(c *gin.Context) {
		proposals := []models.Proposal{}

		fromStr := c.Query("from")
		toStr := c.Query("to")

		if fromStr == "" || toStr == "" {
			fromStr = "-1"
			toStr = "-1"
		}

		from, err := strconv.ParseInt(fromStr, 10, 64)
		to, err := strconv.ParseInt(toStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Please, provide valid 'from' and 'to' fields",
			})
			c.Abort()
			return
		}
		result := database.GlobalDB.Offset(int(from)).Limit(int(to)).Preload("User").Find(&proposals)
		fmt.Println(proposals[0].User)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Server error",
			})
			c.Abort()
			return
		}
		c.JSON(200, proposals)
		return

	}
}

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
