package handler

import (
	"tfg/cmd/db/db_access"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users := db_access.GetUsers(db)
		ctx.JSON(200, users)
	}
}
