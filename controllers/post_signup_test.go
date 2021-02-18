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

func TestMain(m *testing.M) {
	database.Init()
	database.InitRedis()
	DB, _ := database.GlobalDB.DB()
	defer DB.Close()
	defer database.RDB.Redis.Close()
	testLoginUser.HashPassword(testLoginPassword)
	testLoginUser.CreateUserRecord()

	code := m.Run()
	database.GlobalDB.Where("email = ?", testLoginUser.Email).Unscoped().Delete(&models.User{})
	database.GlobalDB.Where("email = ?", testPostUser.Email).Unscoped().Delete(&models.User{})
	os.Exit(code)
}

func TestSignUp(t *testing.T) {
	var actualResult models.User

	user := testPostUser

	payload, err := json.Marshal(&user)
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

	assert.Equal(t, user.Name, actualResult.Name)
	assert.Equal(t, user.Email, actualResult.Email)
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
