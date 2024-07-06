package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hasura/go-graphql-client"
	"github.com/kaleabAlemayehu/foodopia/utility"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Body struct {
	Username string `json:"name"  binding:"required"`
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
type Params struct {
	Body Body `json:"params"`
}
type UserActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Params                 `json:"input"`
}

type UserGraphQLError struct {
	Message string `json:"message"`
}
type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
}

const xHasuraAdminSecret = "x-hasura-admin-secret"

func RegisterNewUser(input Body) (userId int64, err error) {
	// create a client
	adminSecret := os.Getenv("ADMIN_SECRET")
	gqlUrl := os.Getenv("GRAPHQL_URL")
	client := graphql.NewClient(gqlUrl, &http.Client{
		Transport: headerRoundTripper{
			setHeaders: func(req *http.Request) {
				req.Header.Set(xHasuraAdminSecret, adminSecret)
			},
			rt: http.DefaultTransport,
		},
	})
	// get the mutation component and create it
	var m struct {
		InsertUsersOne struct {
			ID       int    `graphql:"id"`
			Username string `graphql:"username"`
			Email    string `graphql:"email"`
		} `graphql:"insert_users_one(object: {email: $email, password_hash: $password_hash, username: $username})"`
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	// create variables for the mutation
	variables := map[string]interface{}{
		"email":         input.Email,
		"password_hash": string(hashedPassword),
		"username":      input.Username,
	}
	// create a user
	err = client.Mutate(context.Background(), &m, variables, graphql.OperationName("InsertUser"))
	if err != nil {
		fmt.Printf("what the heck %v\n", err)
	}

	// get id and return it
	return int64(m.InsertUsersOne.ID), nil

}
func Signup(c *gin.Context) {

	// read the data from the request
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		panic(err)
	}

	//  marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	// parse it to userActionPayload
	var actionPayload UserActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userId, err := RegisterNewUser(actionPayload.Input.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"message": err,
		})
	}

	/*
		jwt.MapClaims{
			"id":          user.ID,
			"firstName":   user.FirstName,
			"lastName":    user.LastName,
			"email":       user.Email,
			"phoneNumber": user.PhoneNumber,
			"exp":         expiryDate,
			"https://hasura.io/jwt/claims": map[string]interface{}{
				"x-hasura-default-role":  "user",
				"x-hasura-allowed-roles": [2]string{"user", "admin"},
				"x-hasura-user-id":       strconv.Itoa(user.ID),
			},

	*/

	// Generate JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"email":    actionPayload.Input.Body.Email,
		"username": actionPayload.Input.Body.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": [2]string{"user", "admin"},
			"x-hasura-user-id":       strconv.Itoa(int(userId)),
		},
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, gin.H{
		"name":  actionPayload.Input.Body.Username,
		"email": actionPayload.Input.Body.Email,
		"token": tokenString,
	})
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

func SendEmail(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"errorCode": 1,
		})
	}
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"errorCode": 2,
		})
	}
	email := gjson.GetBytes(jsonString, "event.data.new.email").String()
	username := gjson.GetBytes(jsonString, "event.data.new.username").String()

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	sender := os.Getenv("GMAIL_USERNAME")
	password := os.Getenv("GMAIL_PASSWORD") // Ideally, use environment variables for security

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
