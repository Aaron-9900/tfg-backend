// models/models_test.go

package models

import (
	"os"
	"testing"
	"tfg/database"

	"github.com/stretchr/testify/assert"
)

var testUser User

func TestMain(m *testing.M) {
	err := database.Init()
	if err != nil {
		panic(err)
	}
	testUser = User{Name: "Test user", Email: "newemail@gmail.com", Password: "123456"}
	testUser.CreateUserRecord()
	m.Run()
	database.GlobalDB.Unscoped().Delete(&testUser)

}

func TestHashPassword(t *testing.T) {
	user := User{
		Password: "secret",
	}

	err := user.HashPassword(user.Password)
	assert.NoError(t, err)

	os.Setenv("passwordHash", user.Password)
}

func TestCreateUserRecord(t *testing.T) {
	var userResult User

	err := database.GlobalDB.AutoMigrate(&User{})
	assert.NoError(t, err)

	user := User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: os.Getenv("passwordHash"),
	}

	err = user.CreateUserRecord()
	assert.NoError(t, err)

	database.GlobalDB.Where("email = ?", user.Email).Find(&userResult)

	database.GlobalDB.Unscoped().Delete(&user)

	assert.Equal(t, "Test User", userResult.Name)
	assert.Equal(t, "test@email.com", userResult.Email)

}

func TestCheckPassword(t *testing.T) {
	hash := os.Getenv("passwordHash")

	user := User{
		Password: hash,
	}

	err := user.CheckPassword("secret")
	assert.NoError(t, err)
}

func TestUserProposals(t *testing.T) {
	proposal := Proposal{UserID: testUser.ID, Name: "Proposal request", Description: "New proposal request", Limit: 1}
	err := proposal.CreateProposalRecord()
	assert.NoError(t, err)

	userProposals, err := testUser.GetUserProposals()
	assert.NoError(t, err)
	assert.NotEmpty(t, userProposals)
	assert.Equal(t, userProposals[0].ID, proposal.ID)

}
