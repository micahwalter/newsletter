package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type rsp struct {
	Email string `json:"email"`
	Stat  string `json:"stat"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	fmt.Printf("Email: %s.\n", request.QueryStringParameters["email"])

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	body := rsp{
		Email: request.QueryStringParameters["email"],
		Stat:  "ok",
	}

	b, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
