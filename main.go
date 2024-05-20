package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kaleabAlemayehu/foodopia/utility"
)

// http://localhost:8080/v1/graphql

func main() {

	err := godotenv.Load(".env")

	utility.CheckError(err, "error to load .env file")

	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + port)

	// client := graphql.NewClient(graphqlEndPoint, nil)
	// Use client...

}
