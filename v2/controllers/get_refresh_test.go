package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/v2/auth"
	"tfg/v2/credentials"
	"tfg/v2/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	response := TokenResponse{}
	_, refreshToken, err := auth.GenerateTokens(testLoginUser.IDString())
	assert.NoError(t, err)
	err = database.RDB.SetUserAndToken(testLoginUser.ID, refreshToken)
	assert.NoError(t, err)

	payload, err := json.Marshal(&testLoginUser)
	assert.NoError(t, err)
	request, err := http.NewRequest("GET", "/api/protected/refresh", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Set("token", refreshToken)
	c.Request = request

	Refresh()(c)
	assert.Equal(t, 200, w.Code)
	_, err = database.RDB.GetUserAndToken(testLoginUser.ID, refreshToken)
	// assert.Error(t, err)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	_, err = database.RDB.GetUserAndToken(testLoginUser.ID, response.RefreshToken)
	assert.NoError(t, err)

	jwtWrapperAccess := auth.JwtWrapper{
		SecretKey: credentials.JwtKey,
		Issuer:    "AuthService",
		ExpirationMs: time.Now().Local().Add(time.Hour*time.Duration(0) +
			time.Minute*time.Duration(10) +
			time.Second*time.Duration(0)).Unix(),
	}

	_, err = jwtWrapperAccess.ValidateToken(response.AccessToken)
	assert.NoError(t, err)
}

func TestIncorrectToken(t *testing.T) {
	refreshToken := "incorrect token"

	payload, err := json.Marshal(&testLoginUser)
	assert.NoError(t, err)
	request, err := http.NewRequest("GET", "/api/protected/refresh", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("id", testLoginUser.IDString())
	c.Set("token", refreshToken)
	c.Request = request

	Refresh()(c)
	assert.Equal(t, 401, w.Code)

}
