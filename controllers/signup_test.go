// controllers/public_test.go

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tfg/database"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testLoginUser = models.User{Email: "logintest@email.com", Name: "test"}
var testLoginPassword = "secret"
var testPostUser = models.User{
	Name:     "Test User",
	Email:    "posttest@email.com",
	Password: "secret",
}
var testProposal = models.Proposal{
	Name:        "Proposal test",
	Description: "Proposal description",
	Limit:       1,
}

func TestMain(m *testing.M) {
	database.TestInit("test_tfg")
	database.InitRedis()
	DB, _ := database.GlobalDB.DB()
	defer DB.Close()
	defer database.RDB.Redis.Close()
	testLoginUser.HashPassword(testLoginPassword)
	testLoginUser.CreateUserRecord()
	testProposal.UserID = testLoginUser.ID
	testProposal.CreateProposalRecord()

	code := m.Run()
	database.GlobalDB.Where("id = ?", testProposal.ID).Unscoped().Delete(&models.Proposal{})
	database.GlobalDB.Where("id = ?", testLoginUser.ID).Unscoped().Delete(&models.User{})
	database.GlobalDB.Where("id = ?", testPostUser.ID).Unscoped().Delete(&models.User{})
	os.Exit(code)
}

func TestSignUp(t *testing.T) {
	var actualResult models.User

	payload, err := json.Marshal(&testPostUser)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	Signup()(c)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &actualResult)
	assert.NoError(t, err)
	testPostUser.ID = actualResult.ID
	assert.Equal(t, testPostUser.Name, actualResult.Name)
	assert.Equal(t, testPostUser.Email, actualResult.Email)
}

func TestSignUpInvalidJSON(t *testing.T) {
	user := "test"

	payload, err := json.Marshal(&user)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	Signup()(c)

	assert.Equal(t, 400, w.Code)
}
