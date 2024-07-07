package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hasura/go-graphql-client"
	"github.com/kaleabAlemayehu/foodopia/models"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
}

const xHasuraAdminSecret = "x-hasura-admin-secret"

func RegisterNewUser(input models.Payload) (userId int64, err error) {
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

	//marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	// parse it to userActionPayload
	var actionPayload models.UserActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userId, err := RegisterNewUser(actionPayload.Input.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"message": err,
		})
	}

	// add claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"email":    actionPayload.Input.Payload.Email,
		"username": actionPayload.Input.Payload.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": [2]string{"user", "admin"},
			"x-hasura-user-id":       strconv.Itoa(int(userId)),
		},
	})
	// Generate JWT token
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, gin.H{
		"name":  actionPayload.Input.Payload.Username,
		"email": actionPayload.Input.Payload.Email,
		"token": tokenString,
	})
}
func CheckUser(input models.Payload) (models.Payload, error) {
	// create client
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
	// create query
	var q struct {
		Users []struct {
			Username     string `graphql:"username"`
			PasswordHash string `graphql:"password_hash"`
			Email        string `graphql:"email"`
			Id           int64  `graphql:"id"`
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	// create varaiable
	variable := map[string]interface{}{
		"email": string(input.Email),
	}

	// send request to hasura
	err := client.Query(context.Background(), &q, variable, graphql.OperationName("CheckUser"))
	if err != nil {
		log.Fatalf("error Occured: %v\n", err)
	}
	// check if the password is correct
	storedPassword := q.Users[0].PasswordHash
	if storedPassword == "" {
		panic("there is no password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)); err != nil {
		panic("password is not the same!")
	}

	return models.Payload{
		Id:       q.Users[0].Id,
		Email:    q.Users[0].Email,
		Username: q.Users[0].Username,
	}, nil

}
func Login(c *gin.Context) {
	// read the data from the request
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		panic(err)
	}

	//marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	// unmashall it to models.UserActionPayload
	var actionPayload models.UserActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	// check if the user is legit
	res, err := CheckUser(actionPayload.Input.Payload)
	if err != nil {
		panic("there is fucking error")
	}

	// create a claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       res.Id,
		"email":    res.Email,
		"username": res.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": [2]string{"user", "admin"},
			"x-hasura-user-id":       strconv.Itoa(int(res.Id)),
		},
	})
	// generate jwt token
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       res.Id,
		"username": res.Username,
		"email":    res.Email,
		"token":    tokenString,
	})

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
