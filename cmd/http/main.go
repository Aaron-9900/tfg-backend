package http

import (
	handler "tfg/cmd/http/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run starts the server
func Run(db *gorm.DB) {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/users", handler.GetUsers(db))
	r.POST("/users", handler.PostUser(db))
	r.GET("/user", handler.GetUser(db))
	r.Run("localhost:3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}