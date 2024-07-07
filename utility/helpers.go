package utility

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaleabAlemayehu/foodopia/models"
)

func CreateToken(user models.Payload) string {
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
		panic("unable to stringify token")
	}

	return tokenString
}
func SaveImageToFile(input models.ImageUploadArgs) (models.SaveImageOutput, error) {
	var image models.SaveImageOutput

	// create a decoder with the base64 string from request
	dec, err := base64.StdEncoding.DecodeString(string(input.Base64Str))
	if err != nil {
		panic(err)
	}

	dir, err := filepath.Abs("./uploads")
	if err != nil {
		panic(err)
	}
	// create file and wait to close it after the function is about to return
	file, err := os.Create(filepath.Join(dir, input.Name))
	if err != nil {
		// panic("unable to create a file in the upload directory!")
		panic(err)
	}
	defer file.Close()
	// write the byte to the file
	if _, err = file.Write(dec); err != nil {
		panic(err)
	}
	//  save the file
	if err := file.Sync(); err != nil {
		panic(err)
	}

	image.ImageUrl, err = filepath.Abs(filepath.Join(dir, input.Name))
	if err != nil {
		panic(err)
	}
	return image, err
}
