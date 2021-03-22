// middlewares/authz.go

package middlewares

import (
	"strings"
	"tfg/auth"
	"tfg/credentials"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(400, "Incorrect Format of Authorization Token")
			c.Abort()
			return
		}

		jwtWrapper := auth.JwtWrapper{
			SecretKey: credentials.JwtKey,
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil && err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(400, err.Error())
			c.Abort()
			return
		}
		c.Set("id", claims.ID)
		c.Set("token", clientToken)
		c.Next()

	}
}
