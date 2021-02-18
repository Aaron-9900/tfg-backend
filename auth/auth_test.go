// auth/auth_test.go

package auth

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	jwtWrapper := JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
		ExpirationMs: time.Now().Local().Add(time.Hour*time.Duration(0) +
			time.Minute*time.Duration(10) +
			time.Second*time.Duration(0)).Unix(),
	}

	generatedToken, err := jwtWrapper.GenerateToken("1234")
	assert.NoError(t, err)

	os.Setenv("testToken", generatedToken)
}

func TestValidateToken(t *testing.T) {
	encodedToken := os.Getenv("testToken")

	jwtWrapper := JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.NoError(t, err)

	assert.Equal(t, "1234", claims.ID)
	assert.Equal(t, "AuthService", claims.Issuer)
}
