package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type rsp struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Stat  string `json:"stat"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// grab the email and confirmation code from the query string params
	email := request.QueryStringParameters["email"]
	code := request.QueryStringParameters["code"]

	// make sure we have a valid value for both params

	// get dynamoDB item for email address

	// return an error if it doesn't exist

	// return an error if the code doesn't match

	// update dynamdoDB confirmed true
	
	// parrot back the email for testing
	body := rsp{
		Email: email,
		Code:  code,
		Stat:  "ok",
	}

	b, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
