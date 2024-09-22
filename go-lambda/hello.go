package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MyEvent struct {
	Name string `json:"name"`
}

type TableItem struct {
	Name      string
	Timestamp string
}

func newTableItem(name string) TableItem {
	return TableItem{
		Name:      name,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received event is nil")
	}

	// we create a new session to the AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// we create a new DynamoDB client
	svc := dynamodb.New(sess)

	// we create a new item to be inserted in the table
	item := newTableItem(event.Name)

	// we marshal the item into a map
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}

	// we construct the input request
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	// we put the item in the table
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return aws.String(fmt.Sprintf("Hello %s!", event.Name)), nil
}

func main() {
	lambda.Start(HandleRequest)
}
