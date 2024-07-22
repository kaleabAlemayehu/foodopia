package utility

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaleabAlemayehu/foodopia/models"
)

func CreateToken(user models.Payload) (string, error) {
	// add claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": [2]string{"user", "admin"},
			"x-hasura-user-id":       strconv.Itoa(int(user.Id)),
		},
	})
	// Generate JWT token
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		return "", err

	}

	return tokenString, nil
}
func SaveImageToFile(input models.ImageUploadArgs) (models.SaveImageOutput, error) {
	var image models.SaveImageOutput

	// create a decoder with the base64 string from request
	dec, err := base64.StdEncoding.DecodeString(string(input.Base64Str))
	if err != nil {
		return models.SaveImageOutput{}, errors.New("unable to create base64 decoder")
	}

	dir, err := filepath.Abs("./upload")
	if err != nil {

		return models.SaveImageOutput{}, errors.New("unable to find upload's directory absolute path")
	}
	// create file and wait to close it after the function is about to return

	file, err := os.Create(filepath.Join(dir, input.FileName))
	if err != nil {

		return models.SaveImageOutput{}, errors.New("unable to create a file in the upload directory")
	}
	defer file.Close()
	// write the byte to the file
	if _, err = file.Write(dec); err != nil {

		return models.SaveImageOutput{}, errors.New("unable to write the byte to the file")
	}
	//  save the file
	if err := file.Sync(); err != nil {

		return models.SaveImageOutput{}, errors.New("unable to save the file")
	}

	image.ImageUrl = fmt.Sprintf(`http://localhost:9000/images/%v`, input.FileName)
	// if err != nil {

	// 	return models.SaveImageOutput{}, errors.New("unable to find the absolute path of a file")
	// }
	return image, err
}
