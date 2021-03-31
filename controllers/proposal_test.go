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

func TestPostProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Name:        "Test Proposal",
		Description: "Test proposal description",
		Limit:       10,
		Rate:        1,
		UserID:      testLoginUser.ID,
	}

	response := models.Proposal{}
	payload, err := json.Marshal(&testProposal)
	request, err := http.NewRequest("POST", "/api/protected/proposal", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Request = request

	PostProposal()(c)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, response.Name, testProposal.Name)

	database.GlobalDB.Where("id = ?", response.ID).Delete(&response)
	assert.Equal(t, 200, w.Code)
}

func TestPostWrongProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Description: "Test proposal description",
		Limit:       10,
		Rate:        1,
		UserID:      testLoginUser.ID,
	}

	payload, err := json.Marshal(&testProposal)
	request, err := http.NewRequest("POST", "/api/protected/proposal", bytes.NewBuffer(payload))
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
		Rate:        1,
		UserID:      testLoginUser.ID,
	}

	payload, err := json.Marshal(&testProposal)
	request, err := http.NewRequest("POST", "/api/protected/proposal", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("id", "1241")
	c.Request = request

	PostProposal()(c)

	database.GlobalDB.Where("id = ?", testProposal.ID).Unscoped().Delete(&testProposal)
	assert.Equal(t, 403, w.Code)
}

func TestGetProposal(t *testing.T) {
	var testProposal = models.Proposal{
		Name:        "Test Proposal",
		Description: "Test proposal description",
		Limit:       10,
		Rate:        1,
		UserID:      testLoginUser.ID,
	}
	testProposal.CreateProposalRecord()
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
	assert.Equal(t, requestResponse.User.ID, testProposal.UserID)
	database.GlobalDB.Where("id = ?", testProposal.ID).Unscoped().Delete(&testProposal)

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

func TestGetProposals(t *testing.T) {
	requestResponse := []models.Proposal{}

	request, err := http.NewRequest("GET", "/api/public/proposals", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = request

	GetProposals()(c)
	err = json.Unmarshal(w.Body.Bytes(), &requestResponse)
	if err != nil {
		panic(err)
	}
	proposalFromResponse := models.Proposal{}
	for i := 0; i < len(requestResponse); i++ {
		if requestResponse[i].ID == testProposal.ID {
			proposalFromResponse = requestResponse[i]
		}
	}
	assert.Equal(t, proposalFromResponse.ID, testProposal.ID)
}

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
