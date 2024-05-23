package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(c *gin.Context) {
	// get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	//  Decode / validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("MY_SECRET")), nil
	})

	if err != nil {
		log.Printf("Internal Error: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check the experation
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user with token sub
		userID, ok := claims["user_id"].(string)
		if !ok {
			// Handle case where user ID is a float64
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			userID = fmt.Sprintf("%.0f", userIDFloat)
		}

		c.Set("userID", userID)

		// Continue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

/*


package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

func CheckAuth(c *gin.Context) {
	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(os.Getenv("MY_SECRET")), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expiration
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Assuming the user ID is stored in the "user_id" field
		userID, ok := claims["user_id"].(string)
		if !ok {
			// Handle case where user ID is a float64
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			userID = fmt.Sprintf("%.0f", userIDFloat)
		}

		// Find the user with the token subject (user ID)
		// For demonstration, we're directly attaching the userID to the context
		// In a real scenario, you might query the user from the database
		c.Set("userID", userID)

		// Continue with the request
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func main() {
	r := gin.Default()
	r.Use(CheckAuth)

	// Example route
	r.GET("/protected", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, " + userID,
		})
	})

	r.Run()
}





*/
