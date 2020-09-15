package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	myqueue "github.com/uu64/two-apps/two-back/lib/interface/sqs"
	"github.com/uu64/two-apps/two-back/lib/repository/rooms"
	"github.com/uu64/two-apps/two-back/lib/repository/users"
)

type request events.APIGatewayWebsocketProxyRequest
type response events.APIGatewayProxyResponse

var sqsSvc *sqs.SQS
var dynamoSvc *dynamodb.DynamoDB
var queueName string = "matching"

func getMessage() ([]*sqs.Message, error) {
	message, err := myqueue.ReceiveMessage(sqsSvc, queueName)
	return message.Messages, err
}

func createRoom(connectionID string) (string, error) {
	var roomID string

	// create room
	roomID, err := rooms.Create(dynamoSvc, connectionID)
	if err != nil {
		return roomID, err
	}

	// send message to sqs and wait a new challenger
	err = myqueue.SendMessage(sqsSvc, queueName, roomID)
	return roomID, err
}

func addUser(connectionID string, roomID string) error {
	return users.Create(dynamoSvc, connectionID, roomID)
}

func updateRoom(roomID string, connectionID string, receiptHandle string) error {
	// update room
	err := rooms.AddUser(dynamoSvc, roomID, connectionID)
	if err != nil {
		return err
	}

	// delete message
	err = myqueue.DeleteMessage(sqsSvc, queueName, receiptHandle)
	return err
}

func handler(ctx context.Context, request request) (response, error) {
	fmt.Println("connected!!!!!!!")

	var err error
	messages, err := getMessage()
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, nil
	}

	var roomID string
	connectionID := request.RequestContext.ConnectionID
	if len(messages) == 0 {
		fmt.Println("create room")
		roomID, err = createRoom(connectionID)
	} else {
		roomID = *messages[0].Body

		fmt.Println("match complete")
		err = updateRoom(roomID, connectionID, *messages[0].ReceiptHandle)
	}
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, nil
	}

	err = addUser(connectionID, roomID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, nil
	}

	return response{StatusCode: 200}, nil
}

func init() {
	session := session.New()
	sqsSvc = sqs.New(session)
	dynamoSvc = dynamodb.New(session)
}

func main() {
	lambda.Start(handler)
}
