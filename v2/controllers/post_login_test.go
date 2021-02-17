package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/v2/database"
	"tfg/v2/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	user := LoginPayload{
		Email:    testLoginUser.Email,
		Password: testLoginPassword,
	}
	dbUser := models.User{}
	tokens := TokenResponse{}
	payload, err := json.Marshal(&user)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request
	database.GlobalDB.AutoMigrate(&models.User{})

	Login()(c)

	err = json.Unmarshal(w.Body.Bytes(), &tokens)
	if err != nil {
		panic(err)
	}
	database.GlobalDB.Unscoped().Where("email = ?", user.Email).First(&dbUser)
	val, err := database.RDB.GetUserAndToken(dbUser.ID, tokens.RefreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, val)
	fmt.Println("key", val)
	assert.Equal(t, 200, w.Code)

}

func TestLoginInvalidJSON(t *testing.T) {
	user := "test"

	payload, err := json.Marshal(&user)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	Login()(c)

	assert.Equal(t, 400, w.Code)
}

func TestLoginInvalidCredentials(t *testing.T) {
	user := LoginPayload{
		Email:    testLoginUser.Email,
		Password: "invalid",
	}

	payload, err := json.Marshal(&user)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	err = database.Init()
	assert.NoError(t, err)

	database.GlobalDB.AutoMigrate(&models.User{})

	Login()(c)

	assert.Equal(t, 401, w.Code)
}
