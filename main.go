package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kaleabAlemayehu/foodopia/handlers"
	"github.com/kaleabAlemayehu/foodopia/utility"
)

// http://localhost:8080/v1/graphql

func main() {

	err := godotenv.Load(".env")

	utility.CheckError(err, "error to load .env file")

	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	r.POST("/signup", handlers.Signup)
	r.Run(":" + port)

}
