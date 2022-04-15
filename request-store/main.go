package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type EventRequest struct {
	RequestID  string `json:"ksuid"`
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	RequestURL string `json:"request_url"`
	Date       string `json:"date"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var eventReq EventRequest

	err := json.Unmarshal([]byte(request.Body), &eventReq)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error: %s", err.Error()),
			StatusCode: 500,
		}, nil
	}

	addDocument(eventReq)

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf(
			"Request %s from %d stored successfully",
			eventReq.RequestID,
			eventReq.UserID,
		),
		StatusCode: 200,
	}, nil
}

// Add document to DynamoDB 'Requests' collection
func addDocument(eventReq EventRequest) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(eventReq)
	if err != nil {
		log.Fatalf("Got error marshalling EventRequest item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Requests"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}

func main() {
	lambda.Start(handler)
}
