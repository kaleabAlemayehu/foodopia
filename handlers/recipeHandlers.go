package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaleabAlemayehu/foodopia/models"
	"github.com/kaleabAlemayehu/foodopia/utility"
)

func Upload(c *gin.Context) {
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		panic(err)
	}
	// fmt.Println(string(jsonData))

	//  marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	// unmarshal itto actionPayload
	var actionPayload models.ActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)

	if err != nil {
		panic(err)
	}
	// send it to save the file
	output, err := utility.SaveImageToFile(actionPayload.Input)
	if err != nil {
		panic(err)
	}
	fmt.Println(output.ImageUrl)
	c.JSON(http.StatusOK, gin.H{
		"image_url": output.ImageUrl,
	})

}
