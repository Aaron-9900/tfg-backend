package handler

import (
	"tfg/cmd/db/db_access"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users := db_access.GetUsers()
		ctx.JSON(200, users)
	}
}
