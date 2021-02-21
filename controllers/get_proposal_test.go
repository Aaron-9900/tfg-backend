package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetProposal(t *testing.T) {
	requestResponse := models.Proposal{}
	request, err := http.NewRequest("GET", "/api/protected/proposal", nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("proposal_id", testProposal.IDString())
	request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testProposal.User.IDString())
	c.Request = request

	GetProposal()(c)

	err = json.Unmarshal(w.Body.Bytes(), &requestResponse)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, requestResponse.ID, testProposal.ID)
	assert.Equal(t, requestResponse.Name, testProposal.Name)
	assert.Equal(t, requestResponse.UserID, testProposal.UserID)

}
func TestGetWrongProposal(t *testing.T) {
	requestResponse := models.Proposal{}
	request, err := http.NewRequest("GET", "/api/protected/proposal", nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("proposal_id", "2131")
	request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testProposal.User.IDString())
	c.Request = request

	GetProposal()(c)

	err = json.Unmarshal(w.Body.Bytes(), &requestResponse)
	assert.NoError(t, err)

	assert.Equal(t, 404, w.Code)

}

func TestGetWrongProposalRequest(t *testing.T) {
	requestResponse := models.Proposal{}
	request, err := http.NewRequest("GET", "/api/protected/proposal", nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testProposal.User.IDString())
	c.Request = request

	GetProposal()(c)

	err = json.Unmarshal(w.Body.Bytes(), &requestResponse)
	assert.NoError(t, err)

	assert.Equal(t, 400, w.Code)

}
