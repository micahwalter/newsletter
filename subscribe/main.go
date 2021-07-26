package main

import (
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	rand.Seed(time.Now().UnixNano())

	sayings := []string{
		"hello",
		"heellloo",
	}

	randomSaying := rand.Intn(len(sayings))

	return events.APIGatewayProxyResponse{
		Body:       string(sayings[randomSaying]),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
