package model

type User struct {
	ID       int
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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
