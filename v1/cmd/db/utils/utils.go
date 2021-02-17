package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type HttpError interface {
	Error() string
	ErrorCode() int
}
type RequestError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (err RequestError) Error() string {
	return err.Message
}
func (err RequestError) ErrorCode() int {
	return err.Code
}

var (
	ServerError = RequestError{Message: "Server error", Code: http.StatusInternalServerError}
	AuthError   = RequestError{Message: "Incorrect email/password", Code: http.StatusUnauthorized}
)

func HashAndSalt(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
