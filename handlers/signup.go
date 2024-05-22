package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var CreateUserQueryStr string = `mutation MyMutation($username: String!, $password_hash: String!, $email: String!) {
	insert_foodopia_users_one(object: {email: $email, password_hash: $password_hash, username: $username}) {
		id
		username
		password_hash
	}
}`

func Signup(c *gin.Context) {

	// Load GQL url form environment
	GQLURL := os.Getenv("GRAPHQL_URI")
	if GQLURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "GRAPHQL_URI not set",
		})
		return
	}

	// the user body username/ password / email
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
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
		"query": CreateUserQueryStr,
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
