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
	Solved       bool
}

func getItem(svc *dynamodb.DynamoDB, id string) (User, error) {
	user := User{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, errors.New("user is not exist")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	return user, err
}

// RoomID returns the room-id of the room the user belongs to
func RoomID(svc *dynamodb.DynamoDB, id string) (string, error) {
	user, err := getItem(svc, id)
	return user.RoomID, err
}

// Solved returns whether the user solved the problem
func Solved(svc *dynamodb.DynamoDB, id string) (bool, error) {
	user, err := getItem(svc, id)
	return user.Solved, err
}

// Create creates a user
func Create(svc *dynamodb.DynamoDB, connectionID string, roomID string) error {
	item := User{
		ConnectionID: connectionID,
		RoomID:       roomID,
		Solved:       false,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(userTableName),
	})
	if err != nil {
		return err
	}

	return nil
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

// SolveProblem updates "Solved" to true
func SolveProblem(svc *dynamodb.DynamoDB, id string, userID string) error {
	_, err := svc.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Solved = :s"),
	})

	if err != nil {
		return err
	}

	return nil
}
