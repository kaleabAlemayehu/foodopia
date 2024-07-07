package utility

import (
	"os"
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
