package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context){
	recipeID := c.PostForm("recipe_id")
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
        return
    }

    files := form.File["images"]
    var imageUrls []string

    for _, file := range files {
        // Save the file locally 
        filePath := filepath.Join("uploads", file.Filename)
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }
        imageUrls = append(imageUrls, filePath)
    }

    c.JSON(http.StatusOK, gin.H{"recipe_id": recipeID, "image_urls": imageUrls})
}