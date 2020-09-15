package rooms

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var roomTableName string = "rooms"

// Room is defintion of the rooms table item
type Room struct {
	RoomID  string
	Status  string
	User1ID string
	User2ID string
}

// RoomStatusWaiting is status of the rooms table item
var RoomStatusWaiting string = "WAITING"

// RoomStatusPlaying is status of the rooms table item
var RoomStatusPlaying string = "PLAYING"

// Users returns the connection-id of the user in the room
func Users(svc *dynamodb.DynamoDB, id string) (string, string, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return "", "", err
	}

	if result.Item == nil {
		return "", "", errors.New("room is not exist")
	}

	item := Room{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", "", err
	}
	return item.User1ID, item.User2ID, nil
}

// Create deletes the room with the id
func Create(svc *dynamodb.DynamoDB, userID string) error {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	roomID := uuidObj.String()
	item := Room{
		RoomID:  roomID,
		Status:  RoomStatusWaiting,
		User1ID: userID,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(roomTableName),
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the room with the id
func Delete(svc *dynamodb.DynamoDB, id string) error {
	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
