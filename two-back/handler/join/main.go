package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type request events.APIGatewayWebsocketProxyRequest
type response events.APIGatewayProxyResponse

type room struct {
	RoomID  string
	User1ID string
	User2ID string
}

type user struct {
	ConnectionID string
	RoomID       string
}

var sqsSvc *sqs.SQS
var dynamoSvc *dynamodb.DynamoDB
var matchingQueueName string = "matching"
var roomTableName string = "rooms"
var userTableName string = "users"

func getMessage() ([]sqs.Message, error) {
	messages := []sqs.Message{}

	// receive message from sqs
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &matchingQueueName,
	})
	if err != nil {
		return messages, err
	}

	msgResult, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            urlResult.QueueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(30),
	})

	if len(msgResult.Messages) > 0 {
		messages = append(messages, *msgResult.Messages[0])
	}
	return messages, nil
}

func createRoom(connectionID string) (string, error) {
	// create room
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	roomID := uuidObj.String()
	roomItem := room{
		RoomID:  roomID,
		User1ID: connectionID,
	}
	roomAv, err := dynamodbattribute.MarshalMap(roomItem)
	if err != nil {
		return "", err
	}
	_, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
		Item:      roomAv,
		TableName: aws.String(roomTableName),
	})
	if err != nil {
		return "", err
	}

	// send message to sqs and wait a new challenger
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &matchingQueueName,
	})
	if err != nil {
		return "", err
	}
	_, err = sqsSvc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  &roomID,
		QueueUrl:     urlResult.QueueUrl,
	})
	if err != nil {
		return "", err
	}

	return roomID, nil
}

func addUser(roomID string, connectionID string) error {
	// add user
	userItem := user{
		ConnectionID: connectionID,
		RoomID:       roomID,
	}
	userAv, err := dynamodbattribute.MarshalMap(userItem)
	if err != nil {
		return err
	}
	_, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
		Item:      userAv,
		TableName: aws.String(userTableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func updateRoom(roomID string, connectionID string, receiptHandle string) error {
	// update room
	_, err := dynamoSvc.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(connectionID),
			},
		},
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(roomID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set User2ID = :id"),
	})
	if err != nil {
		return err
	}

	// delete message
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &matchingQueueName,
	})
	if err != nil {
		return err
	}
	_, err = sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      urlResult.QueueUrl,
		ReceiptHandle: &receiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
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

	err = addUser(roomID, connectionID)
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
