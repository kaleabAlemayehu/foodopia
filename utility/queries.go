package utility

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/kaleabAlemayehu/foodopia/models"
	"golang.org/x/crypto/bcrypt"
)

const xHasuraAdminSecret = "x-hasura-admin-secret"

type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
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

func RegisterNewUser(input models.Payload) (models.Payload, error) {
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
		panic(err)
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
	return models.Payload{
		Id:       int64(m.InsertUsersOne.ID),
		Username: m.InsertUsersOne.Username,
		Email:    m.InsertUsersOne.Email,
	}, nil

}
