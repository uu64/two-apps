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
	Problem []int
}

// RoomStatusWaiting is status of the rooms table item
var RoomStatusWaiting string = "WAITING"

// RoomStatusPlaying is status of the rooms table item
var RoomStatusPlaying string = "PLAYING"

func getItem(svc *dynamodb.DynamoDB, id string) (Room, error) {
	room := Room{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return room, err
	}

	if result.Item == nil {
		return room, errors.New("room is not exist")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &room)
	return room, err
}

// Users returns the connection-id of the user in the room
func Users(svc *dynamodb.DynamoDB, id string) ([]string, error) {
	room, err := getItem(svc, id)
	return []string{room.User1ID, room.User2ID}, err
}

// Status returns the status of the room
func Status(svc *dynamodb.DynamoDB, id string) (string, error) {
	room, err := getItem(svc, id)
	return room.Status, err
}

// Problem returns the problem of the room
func Problem(svc *dynamodb.DynamoDB, id string) ([]int, error) {
	room, err := getItem(svc, id)
	return room.Problem, err
}

// Create creates a room and returns the room-id
func Create(svc *dynamodb.DynamoDB, userID string) (string, error) {
	var roomID string

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return roomID, err
	}

	roomID = uuidObj.String()
	item := Room{
		RoomID:  roomID,
		Status:  RoomStatusWaiting,
		User1ID: userID,
		User2ID: "",
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return roomID, err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(roomTableName),
	})
	if err != nil {
		return roomID, err
	}

	return roomID, nil
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

// AddUser adds the user to the room
func AddUser(svc *dynamodb.DynamoDB, id string, userID string) error {
	_, err := svc.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#st": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(userID),
			},
			":st": {
				S: aws.String(RoomStatusPlaying),
			},
		},
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set User2ID = :id, #st = :st"),
	})

	if err != nil {
		return err
	}

	return nil
}

// SetProblem sets a problem to the room
func SetProblem(svc *dynamodb.DynamoDB, id string, problem []int) error {
	av, err := dynamodbattribute.Marshal(problem)
	if err != nil {
		return err
	}

	_, err = svc.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": av,
		},
		TableName: aws.String(roomTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Problem = :p"),
	})

	return err
}
