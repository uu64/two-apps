package users

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var userTableName string = "users"

// User is defintion of the users table item
type User struct {
	ConnectionID string
	RoomID       string
}

// RoomStatusWaiting is status of the rooms table item
var RoomStatusWaiting string = "WAITING"

// RoomStatusPlaying is status of the rooms table item
var RoomStatusPlaying string = "PLAYING"

// RoomID returns the room-id of the room the user belongs to
func RoomID(svc *dynamodb.DynamoDB, id string) (string, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", errors.New("user is not exist")
	}

	item := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}
	return item.RoomID, nil
}

// Delete deletes the user with the id
func Delete(svc *dynamodb.DynamoDB, id string) error {
	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
