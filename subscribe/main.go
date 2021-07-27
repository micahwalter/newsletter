package main

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type rsp struct {
	Email string `json:"email"`
	Stat  string `json:"stat"`
}

type Subscriber struct {
	Email string
	Confirmed bool
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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	data := Subscriber{
		Email: email,
		Confirmed: false,
	}

	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		exitErrorf("Got error marshalling new item:", av, err)
	}

	tableName := "Subscribers"
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		exitErrorf("Got error calling PutItem:", err)
	}

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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	lambda.Start(Handler)
}
