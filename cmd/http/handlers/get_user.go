package handler

import (
	"tfg/cmd/db/db_access"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("username")
		users := db_access.GetUser(db, name)
		ctx.JSON(200, users)
	}
}
