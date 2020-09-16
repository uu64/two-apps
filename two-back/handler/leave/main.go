package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/uu64/two-apps/two-back/lib/interface/ws"
	"github.com/uu64/two-apps/two-back/lib/repository/rooms"
	"github.com/uu64/two-apps/two-back/lib/repository/users"
)

type request events.APIGatewayWebsocketProxyRequest
type response events.APIGatewayProxyResponse

var dynamoSvc *dynamodb.DynamoDB
var agwSvc *apigatewaymanagementapi.ApiGatewayManagementApi

func getRoomID(connectionID string) (string, error) {
	return users.RoomID(dynamoSvc, connectionID)
}

func getUsers(roomID string) ([]string, error) {
	return rooms.Users(dynamoSvc, roomID)
}

func deleteUser(userID string) error {
	return users.Delete(dynamoSvc, userID)
}

func deleteRoom(roomID string) error {
	return rooms.Delete(dynamoSvc, roomID)
}

func handler(ctx context.Context, request request) (response, error) {
	fmt.Println("disconnected!!!!!!!")

	connectionID := request.RequestContext.ConnectionID
	roomID, err := getRoomID(connectionID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

	userList, err := getUsers(roomID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

	endpoint := fmt.Sprintf("https://%s/%s",
		request.RequestContext.DomainName, request.RequestContext.Stage)
	user1ID := userList[0]
	user2ID := userList[1]
	if user1ID == connectionID {
		ws.Disconnect(agwSvc, endpoint, user2ID)
	}
	if user2ID == connectionID {
		ws.Disconnect(agwSvc, endpoint, user1ID)
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
