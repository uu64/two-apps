package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

var dynamoSvc *dynamodb.DynamoDB
var agwSvc *apigatewaymanagementapi.ApiGatewayManagementApi
var roomTableName string = "rooms"
var userTableName string = "users"

func getRoomID(connectionID string) (string, error) {
	userResult, err := dynamoSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(connectionID),
			},
		},
	})
	if err != nil {
		return "", err
	}

	if userResult.Item == nil {
		return "", errors.New("user is not exist")
	}

	userItem := user{}
	err = dynamodbattribute.UnmarshalMap(userResult.Item, &userItem)
	if err != nil {
		return "", err
	}
	return userItem.RoomID, nil
}

func getUsers(roomID string) (string, string, error) {
	roomResult, err := dynamoSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(roomID),
			},
		},
	})
	if err != nil {
		return "", "", err
	}

	if roomResult.Item == nil {
		return "", "", errors.New("room is not exist")
	}

	roomItem := room{}
	err = dynamodbattribute.UnmarshalMap(roomResult.Item, &roomItem)
	if err != nil {
		return "", "", err
	}
	return roomItem.User1ID, roomItem.User2ID, nil
}

func deleteUser(userID string) error {
	_, err := dynamoSvc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteRoom(roomID string) error {
	_, err := dynamoSvc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(roomID),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteConnection(endpoint string, connectionID string) {
	agwSvc.Endpoint = endpoint
	agwSvc.DeleteConnection(&apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: &connectionID,
	})
}

func handler(ctx context.Context, request request) (response, error) {
	fmt.Println("disconnected!!!!!!!")

	connectionID := request.RequestContext.ConnectionID
	roomID, err := getRoomID(connectionID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, nil
	}

	user1ID, user2ID, err := getUsers(roomID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, nil
	}

	endpoint := fmt.Sprintf("https://%s/%s",
		request.RequestContext.DomainName, request.RequestContext.Stage)
	if user1ID == connectionID {
		deleteConnection(endpoint, user2ID)
	}
	if user2ID == connectionID {
		deleteConnection(endpoint, user1ID)
	}

	deleteUser(user1ID)
	deleteUser(user2ID)
	deleteRoom(roomID)

	return response{StatusCode: 200}, nil
}

func init() {
	session := session.New()
	dynamoSvc = dynamodb.New(session)
	agwSvc = apigatewaymanagementapi.New(session)
}

func main() {
	lambda.Start(handler)
}
