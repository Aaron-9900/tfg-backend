package handler

import (
	"log"
	"tfg/cmd/db"
	"tfg/cmd/db/model"
	"tfg/cmd/http/types"
	"tfg/common/auth"

	"tfg/credentials"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload types.LoginPayload
		var user model.User

		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(400, gin.H{
				"msg": "invalid json",
			})
			ctx.Abort()
			return
		}

		result := db.GlobalDB.Where("email = ?", payload.Email).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(401, gin.H{
				"msg": "invalid user credentials",
			})
			ctx.Abort()
			return
		}
		isCorrect := user.IsCorrectPassword(payload.Password)
		if !isCorrect {
			log.Println(err)
			ctx.JSON(401, gin.H{
				"msg": "invalid user credentials",
			})
			ctx.Abort()
			return
		}
		jwtWrapper := auth.JwtWrapper{
			SecretKey:       credentials.JwtKey,
			Issuer:          "AuthService",
			ExpirationHours: 24,
		}

		signedToken, err := jwtWrapper.GenerateToken(user.Email)
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{
				"msg": "error signing token",
			})
			ctx.Abort()
			return
		}

		tokenResponse := types.LoginResponse{
			Token: signedToken,
		}

		ctx.JSON(200, tokenResponse)

		return
	}
}
