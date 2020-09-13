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
	User1Id string
	User2Id string
}

type user struct {
	ConnectionID string
	RoomID       string
}

func getMessage() ([]sqs.Message, error) {
	messages := []sqs.Message{}

	// create session
	var config *aws.Config
	newSession, err := session.NewSession(config)
	if err != nil {
		return messages, err
	}

	// receive message from sqs
	sqsSvc := sqs.New(newSession)
	queueName := "matching"
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
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
	// create session
	var config *aws.Config
	newSession, err := session.NewSession(config)
	if err != nil {
		return "", err
	}

	// create room
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	dynamoSvc := dynamodb.New(newSession)

	roomID := uuidObj.String()
	roomItem := room{
		RoomID:  roomID,
		User1Id: connectionID,
	}
	roomAv, err := dynamodbattribute.MarshalMap(roomItem)
	if err != nil {
		return "", err
	}
	_, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
		Item:      roomAv,
		TableName: aws.String("rooms"),
	})
	if err != nil {
		return "", err
	}

	// send message to sqs and wait a new challenger
	sqsSvc := sqs.New(newSession)
	queueName := "matching"
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
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
	// create session
	var config *aws.Config
	newSession, err := session.NewSession(config)
	if err != nil {
		return err
	}

	// add user
	dynamoSvc := dynamodb.New(newSession)
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
		TableName: aws.String("users"),
	})
	if err != nil {
		return err
	}

	return nil
}

func updateRoom(roomID string, connectionID string, receiptHandle string) error {
	// create session
	var config *aws.Config
	newSession, err := session.NewSession(config)
	if err != nil {
		return err
	}

	// update room
	dynamoSvc := dynamodb.New(newSession)
	_, err = dynamoSvc.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(connectionID),
			},
		},
		TableName: aws.String("rooms"),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(roomID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set User2Id = :id"),
	})
	if err != nil {
		return err
	}

	// delete message
	sqsSvc := sqs.New(newSession)
	queueName := "matching"
	urlResult, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
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

func Handler(ctx context.Context, request request) (response, error) {
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

func main() {
	lambda.Start(Handler)
}
