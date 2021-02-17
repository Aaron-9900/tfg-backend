// main.go

package main

import (
	"log"
	"tfg/v2/controllers"
	"tfg/v2/database"
	"tfg/v2/middlewares"
	"tfg/v2/models"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	api := r.Group("/api")
	{
		public := api.Group("/public")
		{
			public.POST("/login", controllers.Login())
			public.POST("/signup", controllers.Signup())
		}

		// here
		protected := api.Group("/protected").Use(middlewares.Authz())
		{
			protected.GET("/profile", controllers.Profile())
			protected.GET("/refresh", controllers.Refresh())
		}
	}

	return r
}

func main() {
	err := database.Init()
	DB, _ := database.GlobalDB.DB()
	defer DB.Close()
	if err != nil {
		log.Fatalln("could not create database", err)
	}
	err = database.InitRedis()
	defer database.RDB.Redis.Close()
	if err != nil {
		log.Fatalln("could not create database", err)
	}

	database.GlobalDB.AutoMigrate(&models.User{})

	r := setupRouter()
	r.Run(":3000")
}
