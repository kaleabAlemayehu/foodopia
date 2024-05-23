package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaleabAlemayehu/foodopia/utility"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Body struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	query := utility.CreateUserQueryStr

	// Load GQL url form environment
	var GQLURL string = os.Getenv("GRAPHQL_URI")
	if GQLURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "GRAPHQL_URI not set",
		})
		return
	}
	var body Body
	// the user body username/ password / email
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to hash password",
		})
		return
	}
	log.Printf("username: %s\n password: %s\n email: %s\n", body.Username, body.Password, body.Email)
	// create the user
	payload := map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"username":      body.Username,
			"password_hash": string(hashedPassword),
			"email":         body.Email,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to marshal payload",
		})
		return
	}

	// Perform the HTTP request
	res, err := http.Post(GQLURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}
	defer res.Body.Close()

	// Read the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read response",
		})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	// get query form query.go
	query := utility.CheckUser
	// loading GQL URL
	var GQLURL string = os.Getenv("GRAPHQL_URI")
	if GQLURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "GRAPHQL_URI not set",
		})
		return
	}
	// get email and pass of from req body
	var body Body
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}
	// preparing payload
	checkPayload := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"email": body.Email,
		},
	}

	checkPayloadBytes, err := json.Marshal(checkPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to marshal check payload",
		})
		return
	}

	//  lookup registered user
	checkRes, err := http.Post(GQLURL, "application/json", bytes.NewBuffer(checkPayloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check user",
		})
		return
	}
	defer checkRes.Body.Close()
	var checkResult map[string]interface{}
	if err := json.NewDecoder(checkRes.Body).Decode(&checkResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read check response",
		})
		return
	}

	users, ok := checkResult["data"].(map[string]interface{})["foodopia_users"].([]interface{})
	if !ok || len(users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	user := users[0].(map[string]interface{})
	storedPasswordHash := user["password_hash"].(string)

	// Compare provided password with stored password hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user["id"],
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	//
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	// Send it back
	c.JSON(http.StatusOK, gin.H{
		"message": "succefully logged in",
	})

}

func GeneralFunc(c *gin.Context) {

}
