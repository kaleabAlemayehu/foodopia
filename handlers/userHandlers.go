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

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CreateUserOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type SignupResponse struct {
	Data struct {
		InsertUsersOne User `json:"insert_users_one"`
	} `json:"data"`
}

const xHasuraAdminSecret = "x-hasura-admin-secret"

func RegisterNewUser(input Body) (userId string, err error) {
	// create a client
	adminSecret := os.Getenv("ADMIN_SECRET")
	client := graphql.NewClient(serverEndpoint, &http.Client{
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
		} `graphql:"insert_users_one(object: {email: $email, password_hash: $password_hash, username: $username})"`
	}

	// create variables for the mutation

	// create a user

	// get the returns and get id

	// retrun the id

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

	res, err := RegisterNewUser(actionPayload.Input.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"message": err,
		})
	}

	fmt.Printf("res.Body: %v\n", res.Body)
	defer res.Body.Close()
	// unmarshaling to json
	var resByte map[string]interface{}

	resJson, err := json.Marshal(resByte)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("json %v", string(resJson))

	// Generate JWT token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"email": res
	// 	"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	// })
	// tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "failed to generate token",
	// 	})
	// 	return
	// }
	// // Read the response
	// var result map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "failed to read response",
	// 	})
	// 	return
	// }

	// // Respond with the result
	// c.JSON(http.StatusOK, result)
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
