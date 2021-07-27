package main

import (
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type rsp struct {
	Email string `json:"email"`
	Stat  string `json:"stat"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// grab the email from the query string params
	email := request.QueryStringParameters["email"]

	// validate that its a real email address
	fmt.Printf("%18s valid: %t\n", email, valid(email))

	if !valid(email) {
		body := rsp{
			Email: email,
			Stat:  "error",
		}

		b, _ := json.Marshal(body)

		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	}

	// attempt to insert the email into the database

	// parrot back the email for testing
	body := rsp{
		Email: email,
		Stat:  "ok",
	}

	b, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func main() {
	lambda.Start(Handler)
}
