package handler

import (
	"tfg/cmd/db/db_access"

	"github.com/gin-gonic/gin"
)

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("username")
		users := db_access.GetUser(name)
		ctx.JSON(200, users)
	}
}
