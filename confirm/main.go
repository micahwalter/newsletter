package main

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type rsp struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Stat  string `json:"stat"`
}

type errorMessage struct {
	Stat    string `json:"stat"`
	Message string `json:"message"`
}

type Subscriber struct {
	Email            string
	Confirmed        bool
	ConfirmationCode int
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// grab the email and confirmation code from the query string params
	email := request.QueryStringParameters["email"]
	code := request.QueryStringParameters["code"]

	// make sure we have a valid email

	if !emailValid(email) {
		body := errorMessage{
			Stat:    "error",
			Message: fmt.Sprint("Invalid email address: %v", email),
		}

		b, _ := json.Marshal(body)
		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	}

	// get dynamoDB item for email address

	tableName := "Subscribers"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	})

	// in case we get an error from DynamoDB
	if err != nil {
		body := errorMessage{
			Stat:    "error",
			Message: fmt.Sprintf("DynamoDB returned an error: %v", err),
		}

		b, _ := json.Marshal(body)
		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	}

	// return an error a record doesn't exist
	if result.Item == nil {
		body := errorMessage{
			Stat:    "error",
			Message: fmt.Sprintf("Failed to find record for: %v", email),
		}

		b, _ := json.Marshal(body)
		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	}

	// unmarshal record
	subscriber := Subscriber{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &subscriber)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record: %v\n", err))
	}

	// return an error if the code doesn't match
	intCode, _ := strconv.Atoi(code)
	if subscriber.ConfirmationCode != intCode {
		body := errorMessage{
			Stat:    "error",
			Message: fmt.Sprintf("Invalid confirmation code: %v", code),
		}

		b, _ := json.Marshal(body)
		return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
	}

	// update dynamdoDB confirmed true
	data := Subscriber{
		Email:     email,
		Confirmed: true,
		ConfirmationCode: intCode,
	}

	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		exitErrorf("Got error marshalling new item:", av, err)
	}

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
		Code:  code,
		Stat:  "ok",
	}

	b, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{Body: string(b), StatusCode: 200}, nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func emailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func main() {
	lambda.Start(Handler)
}
