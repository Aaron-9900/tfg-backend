package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/auth"
	"tfg/database"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {

	_, refreshToken, err := auth.GenerateTokens(testLoginUser.IDString())
	err = database.RDB.SetUserAndToken(testLoginUser.ID, refreshToken)
	assert.NoError(t, err)
	request, err := http.NewRequest("POST", "/api/protected/logout", bytes.NewBuffer(nil))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("token", refreshToken)
	c.Set("id", testLoginUser.IDString())
	c.Request = request
	Logout()(c)

	if err != nil {
		panic(err)
	}
	val, err := database.RDB.GetUserAndToken(testLoginUser.ID, refreshToken)
	assert.Equal(t, redis.Nil, err)
	assert.Equal(t, "", val)
}
