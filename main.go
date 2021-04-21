// main.go

package main

import (
	"log"
	"math/rand"
	"os"
	"tfg/aws"
	"tfg/controllers"
	"tfg/database"
	"tfg/middlewares"
	"tfg/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	awsSession := aws.Init()
	api := r.Group("/api")
	{
		public := api.Group("/public")
		{
			public.POST("/login", controllers.Login())
			public.POST("/signup", controllers.Signup())
			public.GET("/proposals", controllers.GetProposals())
			public.GET("/proposal-types", controllers.GetProposalTypes())
			public.GET("/privacy-templates", controllers.GetPrivacyTemplates())

		}

		// here
		protected := api.Group("/protected").Use(middlewares.Authz())
		{
			protected.GET("/profile", controllers.Profile())
			protected.GET("/refresh", controllers.Refresh())
			protected.GET("/proposal", controllers.GetProposal())
			protected.POST("/logout", controllers.Logout())
			protected.POST("/proposal", controllers.PostProposal())
			protected.GET("/signed-url", controllers.GetProposalSignedUpload(awsSession))
			protected.POST("/submission", controllers.PostProposalSubmission())
			protected.GET("/submission-file", controllers.GetProposalSignedDownload(awsSession))
			protected.GET("/submissions", controllers.GetUserSubmissions())
			protected.POST("/submission-status", controllers.PostSubmissionStatus())
			protected.GET("/user", controllers.GetUserDetail())
			protected.PUT("/user/settings", controllers.PutUserSettings())
		}
	}

	return r
}
func mockData() {
	gofakeit.Seed(0)

	userData := []models.User{}
	proposalData := []models.Proposal{}
	for i := 0; i < 30; i++ {
		user := models.User{
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: "123456"}
		// TODO: user.HashPassword(user.Password) -> Hash when everything is working
		userData = append(userData, user)
	}
	database.GlobalDB.Create(&userData)

	for i := 0; i < 50; i++ {
		uid := rand.Intn(len(userData)) + int(userData[0].ID)
		proposal := models.Proposal{
			UserID:      uint(uid),
			Name:        gofakeit.ProgrammingLanguage(),
			Limit:       gofakeit.Number(0, 1000),
			Description: gofakeit.Paragraph(1, 1, 300, ""),
			Rate:        gofakeit.Float32Range(0.1, 10),
			Type:        "Facebook"}
		proposalData = append(proposalData, proposal)
	}
	database.GlobalDB.Create(&proposalData)
}

func main() {
	err := database.Init()
	DB, _ := database.GlobalDB.DB()
	defer DB.Close()
	if len(os.Args) > 1 && os.Args[1] == "withMocks" {
		mockData()
		return
	}
	if err != nil {
		log.Fatalln("could not create database", err)
	}
	err = database.InitRedis()
	defer database.RDB.Redis.Close()
	if err != nil {
		log.Fatalln("could not create database", err)
	}

	r := setupRouter()
	r.Run(":3000")
}
