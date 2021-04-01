package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostSubmission(t *testing.T) {
	postSubmissionRequest := &postSubmissionRequest{
		UserID:     testLoginUser.ID,
		ProposalID: testProposal.ID,
		FileName:   "test_name",
	}
	response := models.Submission{}
	payload, err := json.Marshal(&postSubmissionRequest)
	request, err := http.NewRequest("POST", "/api/protected/submission", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Request = request
	PostProposalSubmission()(c)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	database.GlobalDB.Where("id = ?", response.ID).Delete(&response)
	assert.Equal(t, 200, w.Code)

}

func TestForbiddenSubmission(t *testing.T) {
	postSubmissionRequest := &postSubmissionRequest{
		UserID:     testLoginUser.ID,
		ProposalID: testProposal.ID,
		FileName:   "test_name",
	}
	response := models.Submission{}
	payload, err := json.Marshal(&postSubmissionRequest)
	request, err := http.NewRequest("POST", "/api/protected/submission", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", "1000")
	c.Request = request
	PostProposalSubmission()(c)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 403, w.Code)
}
