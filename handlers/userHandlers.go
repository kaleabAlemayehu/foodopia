package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/kaleabAlemayehu/foodopia/models"
	"github.com/kaleabAlemayehu/foodopia/utility"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {

	// read the data from the request
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    "BadRequest: invalid user action payload",
		})
		return
	}

	//marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    "BadRequest: invalid user action payload",
		})
		return
	}
	// parse it to userActionPayload
	var actionPayload models.UserActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    "BadRequest: invalid user action payload",
		})
		return
	}

	actionPayload.Input.Payload, err = utility.RegisterNewUser(actionPayload.Input.Payload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    "The User is already registered",
		})
		return
	}

	// Generate JWT token
	tokenString, err := utility.CreateToken(actionPayload.Input.Payload)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    "Internal Server Error: unable to create token",
		})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, gin.H{
		"id":       actionPayload.Input.Payload.Id,
		"email":    actionPayload.Input.Payload.Email,
		"username": actionPayload.Input.Payload.Username,
		"token":    string(tokenString),
	})
}

func Login(c *gin.Context) {
	// read the data from the request
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    err,
		})
	}

	//marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    err,
		})
	}

	// unmashall it to models.UserActionPayload
	var actionPayload models.UserActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    err,
		})
	}
	// check if the user is legit
	res, err := utility.CheckUser(actionPayload.Input.Payload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    err,
		})
	}

	// generate jwt token
	tokenString, err := utility.CreateToken(res)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":       0,
			"email":    "",
			"username": "",
			"token":    "",
			"error":    err,
		})
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
