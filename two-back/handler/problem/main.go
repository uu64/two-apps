package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

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

type data struct {
	Level int `json:"level"`
}

type reply struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func getRoomStatus(connectionID string) (string, error) {
	roomID, err := users.RoomID(dynamoSvc, connectionID)
	if err != nil {
		return roomID, err
	}
	return rooms.Status(dynamoSvc, roomID)
}

func getRoomUsers(connectionID string) ([]string, error) {
	var userList []string

	roomID, err := users.RoomID(dynamoSvc, connectionID)
	if err != nil {
		return userList, err
	}
	return rooms.Users(dynamoSvc, roomID)
}

func onWaiting(endpoint string, connectionID string) error {
	reply := reply{
		Message: "PLEASE_WAIT",
	}

	data, err := json.Marshal(&reply)
	if err != nil {
		return err
	}
	ws.Send(agwSvc, endpoint, []string{connectionID}, data)
	return nil
}

func onPlaying(endpoint string, connectionID string, level int) error {
	problem, err := createProblem(level)
	if err != nil {
		return err
	}

	connectionIDs, err := getRoomUsers(connectionID)
	if err != nil {
		return err
	}

	reply := reply{
		Message: "START_GAME",
		Data:    problem,
	}
	data, err := json.Marshal(&reply)
	if err != nil {
		return err
	}

	ws.Send(agwSvc, endpoint, connectionIDs, data)
	return nil
}

func createProblem(num int) ([]int, error) {
	terms := make([]int, num)

	if num > 10 || num < 0 {
		return terms, errors.New("INVALID_PARAMETER")
	}

	rand.Seed(time.Now().UnixNano())

	sum := 2
	for i := 0; i < num-1; i++ {
		term := rand.Intn(10)
		switch rand.Intn(2) {
		case 0:
			sum = sum + term
		case 1:
			sum = sum - term
		}
		terms[num-1-i] = term
	}
	terms[0] = sum

	return terms, nil
}

func handler(ctx context.Context, request request) (response, error) {
	connectionID := request.RequestContext.ConnectionID
	endpoint := fmt.Sprintf("https://%s/%s",
		request.RequestContext.DomainName, request.RequestContext.Stage)

	// check room status
	status, err := getRoomStatus(connectionID)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

	if status == rooms.RoomStatusWaiting {
		err = onWaiting(endpoint, connectionID)
	}

	if status == rooms.RoomStatusPlaying {
		var data data
		err = json.Unmarshal([]byte(request.Body), &data)
		if err != nil {
			fmt.Println(err)
			return response{StatusCode: 500}, err
		}

		err = onPlaying(endpoint, connectionID, data.Level)
	}

	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

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
