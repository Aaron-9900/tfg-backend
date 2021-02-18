// controllers/protected_test.go

package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"tfg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	var profile models.User

	user := testLoginUser
	request, err := http.NewRequest("GET", "/api/protected/profile", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	c.Set("id", user.IDString())

	Profile()(c)

	err = json.Unmarshal(w.Body.Bytes(), &profile)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)

	log.Println(profile)

	assert.Equal(t, user.Email, profile.Email)
	assert.Equal(t, user.Name, profile.Name)
}

func TestProfileNotFound(t *testing.T) {
	var profile models.User

	request, err := http.NewRequest("GET", "/api/protected/profile", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = request

	c.Set("id", "0")

	Profile()(c)
	err = json.Unmarshal(w.Body.Bytes(), &profile)
	assert.NoError(t, err)
	assert.Equal(t, 404, w.Code)
}
