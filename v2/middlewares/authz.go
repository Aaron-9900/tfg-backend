// middlewares/authz.go

package middlewares

import (
	"fmt"
	"strings"
	"tfg/v2/auth"
	"tfg/v2/credentials"

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
		if err != nil {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}
		fmt.Println(claims)
		c.Set("id", claims.ID)
		c.Set("token", clientToken)
		c.Next()

	}
}
