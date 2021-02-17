package handler

import (
	"net/http"
	"tfg/cmd/db/db_access"
	"tfg/cmd/db/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userPost struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func PostUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := userPost{}
		ctx.Bind(&user)
		validate := validator.New()
		err := validate.Struct(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		finalUser, err := db_access.PostUser(model.User{ID: 0, Name: user.Name, Email: user.Email, Password: user.Password})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, finalUser)
	}
}
