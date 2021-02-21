package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Name:        "Test Proposal",
		Description: "Test proposal description",
		Limit:       10,
		UserID:      testLoginUser.ID,
	}

	request, err := http.NewRequest("POST", "/api/protected/proposal", nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("name", testProposal.Name)
	q.Add("description", testProposal.Description)
	q.Add("limit", fmt.Sprint(testProposal.Limit))
	request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Request = request

	PostProposal()(c)

	database.GlobalDB.Where("id = ?", testProposal.ID).Delete(&testProposal)
	assert.Equal(t, 200, w.Code)
}

func TestPostWrongProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Name:        "Test Proposal",
		Description: "Test proposal description",
		Limit:       10,
		UserID:      testLoginUser.ID,
	}

	request, err := http.NewRequest("POST", "/api/protected/proposal", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Request = request

	PostProposal()(c)

	database.GlobalDB.Where("id = ?", testProposal.ID).Delete(&testProposal)
	assert.Equal(t, 400, w.Code)
}

func TestPostUnAuthProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Name:        "Test Proposal",
		Description: "Test proposal description",
		Limit:       10,
		UserID:      testLoginUser.ID,
	}

	request, err := http.NewRequest("POST", "/api/protected/proposal", nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("name", testProposal.Name)
	q.Add("description", testProposal.Description)
	q.Add("limit", fmt.Sprint(testProposal.Limit))
	request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", "1241")
	c.Request = request

	PostProposal()(c)

	database.GlobalDB.Where("id = ?", testProposal.ID).Unscoped().Delete(&testProposal)
	assert.Equal(t, 500, w.Code)
}
