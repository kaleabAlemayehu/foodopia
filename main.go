package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kaleabAlemayehu/foodopia/handlers"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		panic("error to load .env file")
	}

	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	r.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"POST", "GET", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/upload", handlers.Upload)
	r.POST("/welcome", handlers.SendEmail)
	r.Run(":" + port)
}
