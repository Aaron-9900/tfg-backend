// auth/auth.go

package auth

import (
	"errors"
	"fmt"
	"tfg/credentials"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey    string
	Issuer       string
	ExpirationMs int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	ID string
	jwt.StandardClaims
}

func GenerateTokens(id string) (string, string, error) {
	jwtWrapperAccess := JwtWrapper{
		SecretKey: credentials.JwtKey,
		Issuer:    "AuthService",
		ExpirationMs: time.Now().Local().Add(time.Hour*time.Duration(0) +
			time.Minute*time.Duration(10) +
			time.Second*time.Duration(0)).Unix(),
	}
	jwtWrapperRefresh := JwtWrapper{
		SecretKey: credentials.JwtKey,
		Issuer:    "AuthService",
		ExpirationMs: time.Now().Local().Add(time.Hour*time.Duration(72) +
			time.Minute*time.Duration(0) +
			time.Second*time.Duration(0)).Unix(),
	}

	signedToken, err := jwtWrapperAccess.GenerateToken(fmt.Sprint(id))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwtWrapperRefresh.GenerateToken(fmt.Sprint(id))
	if err != nil {
		return "", "", err
	}
	return signedToken, refreshToken, nil
}

// GenerateToken generates a jwt token
func (j *JwtWrapper) GenerateToken(id string) (signedToken string, err error) {
	claims := &JwtClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: j.ExpirationMs,
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

//ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
