package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type imageUploadArgs struct {
	Name string `json:"name"`
	Base64Str string `json:"base64Str"`
}
type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input imageUploadArgs `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}
type saveImageOutput struct{
	ImageUrl string `json:"image_url"`
}


func saveImageToFile(input imageUploadArgs) (saveImageOutput , error) {
	var image saveImageOutput 
	// fmt.Printf("name: %v", string(input.Name))
	// fmt.Printf("data: %v", string(input.Base64Str))
	// saveDir, err := filepath.Abs("../uploads")
	// if err != nil {
	// 	panic("unable to find save directory")
	// }
	// create a decoder with the base64 string from request
	dec, err := base64.StdEncoding.DecodeString(string(input.Base64Str))
	if err != nil {
		return image, err
	}

	dir, err := filepath.Abs("./uploads")
	if err != nil {
		return image, err
	}
	// os.Create(filepath.Join(dir, filepath.Base(file.Filename)))
	// create file and wait to close it after the function is about to return
	file, err := os.Create(filepath.Join(dir, input.Name))
	if err != nil {
		// panic("unable to create a file in the upload directory!")
		return image, err
	}
	defer file.Close()
	// write the byte to the file
	if _ , err = file.Write(dec); err != nil{
		return image, err
	}
	//  save the file 
	if err := file.Sync(); err != nil {
		return image, err
    }
	
	image.ImageUrl , err = filepath.Abs(filepath.Join(dir, input.Name))
	if err != nil {
		return image, err
	}
	return image , err
}

func Upload (c *gin.Context){
	var jsonData map[string]interface{}
	// convert the byte( i think it is what it is) from gin.context to json
	if err := c.ShouldBindJSON(&jsonData); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"errorCode": 1,
		})
	}
	// fmt.Println(string(jsonData))

	//  marshal it to jsonString (map[string]Interface to json string but bytes)
	jsonString, err:= json.Marshal(jsonData)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"errorCode": 2,
		})
	}
	// unmarshal itto actionPayload
	var actionPayload ActionPayload
	err = json.Unmarshal(jsonString, &actionPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"errorCode": 3,
		})
	}
	// send it to save the file
	output, err := saveImageToFile(actionPayload.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to save file",
			"error": err.Error(),
			"errorCode": 3,
		})
	}
	data , _ := json.Marshal(output)
	fmt.Println(string(data))

	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

// func Uploads(c *gin.Context){
// 	recipeID := c.PostForm("recipe_id")
//     form, err := c.MultipartForm()
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
//         return
//     }
//     files := form.File["images"]
//     var imageUrls []string

//     for _, file := range files {
//         // Save the file locally 
//         filePath := filepath.Join("uploads", file.Filename)
//         if err := c.SaveUploadedFile(file, filePath); err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
//             return
//         }
//         imageUrls = append(imageUrls, filePath)
// 	fmt.Print(filePath)
//     } 
	
//     c.JSON(http.StatusOK, gin.H{"recipe_id": recipeID, "image_urls": imageUrls[0]})
// }
// func Upload( c *gin.Context){
// 	image, err := c.FormFile("image")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{ "error": "Invalid form data", "message": err.Error()})
// 		return
// 	}
// 	dst := "./uploads/" + image.Filename
// 	if err = c.SaveUploadedFile(image, dst); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{ "error" : "unable to upload image", "message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{ "message": "uploaded successfully", "data": image.Filename})
// }

						


