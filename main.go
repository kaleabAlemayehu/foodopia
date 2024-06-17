package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kaleabAlemayehu/foodopia/handlers"
	"github.com/kaleabAlemayehu/foodopia/middlewares"
	"github.com/kaleabAlemayehu/foodopia/utility"
)

func main() {

	err := godotenv.Load(".env")


	utility.CheckError(err, "error to load .env file")

	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	r.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"POST", "GET", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge: 12 * time.Hour,
	}))
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/graphql", middlewares.CheckAuth, handlers.ProxyToHasura)
	r.POST("/upload", handlers.Upload);
	r.POST("/welcome", handlers.SendEmail);
	r.Run(":" + port)
}
