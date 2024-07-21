package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaleabAlemayehu/foodopia/models"
	"github.com/kaleabAlemayehu/foodopia/utility"
)

func Upload(c *gin.Context) {
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"imageUrl": "",
			"error":    "Bad Request: Invalid upload action Payload!",
		})
	}
	// fmt.Println(string(jsonData))

	//  marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"imageUrl": "",
			"error":    "Bad Request: Invalid upload action Payload!",
		})
	}
	// unmarshal itto actionPayload
	var actionPayload models.ActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"imageUrl": "",
			"error":    "Bad Request: Invalid upload action Payload!",
		})
	}
	// send it to save the file
	output, err := utility.SaveImageToFile(actionPayload.Input)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"imageUrl": "",
			"error":    fmt.Sprintf("InternalServerError: unable to upload the file! %v", err.Error()),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"imageUrl": output.ImageUrl,
		"error":    "",
	})

}

func ServeImage(c *gin.Context) {
	// Extract and clean the file path from the request URL
	path := strings.TrimPrefix(c.Request.URL.Path, "/images/")
	fmt.Printf("Requested path: %v\n", path)

	filePath := filepath.Join("./upload", path)
	fmt.Printf("Serving file from path: %v\n", filePath)

	if _, err := filepath.Abs(filePath); err != nil {
		fmt.Printf("File not found: %v\n", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
