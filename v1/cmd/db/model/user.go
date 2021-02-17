package model

import (
	"tfg/cmd/db/utils"
)

type User struct {
	ID           uint
	Name         string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	RefreshToken string `json:"token"`
	TokenExpire  int64  `json:"token_expire"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	DeletedAt    int64  `json:"deleted_at"`
}

type PublicUser struct {
	Name  string `json:"username"`
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (u *User) ToPublicModel() PublicUser {
	return PublicUser{
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *User) IsCorrectPassword(password string) bool {
	return utils.ComparePasswords(u.Password, password)
}

func (u *User) GenerateSecurePassword() string {
	return utils.HashAndSalt(u.Password)
}
