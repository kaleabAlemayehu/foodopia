package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaleabAlemayehu/foodopia/utility"
	"github.com/tidwall/gjson"

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

	users, ok := checkResult["data"].(map[string]interface{})["users"].([]interface{})
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
	log.Printf("user log: %v", user)
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user["id"],
		"exp":time.Now().Add(time.Hour * 24 * 30).Unix(),
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

func ProxyToHasura(c *gin.Context) {
	// Get the Hasura GraphQL endpoint
	hasuraURL := os.Getenv("GRAPHQL_URI")
	if hasuraURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "HASURA_GRAPHQL_URL not set",
		})
		return
	}

	// Read the body from the original request
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Create a new request to forward to Hasura
	req, err := http.NewRequest(c.Request.Method, hasuraURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create request",
		})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send request to Hasura",
		})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read response from Hasura",
		})
		return
	}

	// Copy the response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Write the status code
	c.Writer.WriteHeader(resp.StatusCode)

	// Write the response body
	c.Writer.Write(respBody)
}

func SendEmail(c *gin.Context){
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"errorCode": 1,
		})
	}
	jsonString, err:= json.Marshal(jsonData)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"errorCode": 2,
		})
	}
	email := gjson.GetBytes(jsonString, "event.data.new.email").String()
	username := gjson.GetBytes(jsonString, "event.data.new.username").String()

	smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    sender := os.Getenv("GMAIL_USERNAME")
    password := os.Getenv("GMAIL_PASSWORD")  // Ideally, use environment variables for security

    // Message.
    subject := "Welcome Foodopia\n"
    body := "Dear %v .\nWe Are very Delighted To Have You On Our Platform.\n We hope you will enjoy the great foodie community we have"

    message := []byte(subject + "\n" + fmt.Sprintf(body, username))

    // Authentication.
    auth := smtp.PlainAuth("", sender, password, smtpHost)
/*
*/
    // Sending email.
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{email}, message)
    if err != nil {
        fmt.Println("Error sending email:", err)
        return
    }
    fmt.Println("Email sent successfully!")
}