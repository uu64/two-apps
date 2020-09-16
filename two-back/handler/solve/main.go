package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uu64/two-apps/two-back/lib/interface/ws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/uu64/two-apps/two-back/lib/repository/rooms"
	"github.com/uu64/two-apps/two-back/lib/repository/users"
)

type request events.APIGatewayWebsocketProxyRequest
type response events.APIGatewayProxyResponse

var dynamoSvc *dynamodb.DynamoDB
var agwSvc *apigatewaymanagementapi.ApiGatewayManagementApi

type incoming struct {
	Problem []int    `json:"problem"`
	Answer  []string `json:"answer"`
}

type outgoing struct {
	Message string `json:"message"`
}

func getRoomStatus(connectionID string) (string, error) {
	roomID, err := users.RoomID(dynamoSvc, connectionID)
	if err != nil {
		return roomID, err
	}
	return rooms.Status(dynamoSvc, roomID)
}

func checkAnswer(problem []int, answer []string) bool {
	if len(problem) != len(answer)+1 {
		return false
	}

	num := problem[0]
	for i, v := range answer {
		if v == "p" {
			num = num + problem[i+1]
		} else {
			num = num - problem[i+1]
		}
	}
	if num != 2 {
		return false
	}

	return true
}

func checkChallenger(connectionID string) (bool, error) {
	var solved bool

	roomID, err := users.RoomID(dynamoSvc, connectionID)
	if err != nil {
		return solved, err
	}

	roomUsers, err := rooms.Users(dynamoSvc, roomID)
	if roomUsers[0] == connectionID {
		solved, err = users.Solved(dynamoSvc, roomUsers[1])
	} else if roomUsers[1] == connectionID {
		solved, err = users.Solved(dynamoSvc, roomUsers[0])
	} else {
		err = errors.New("USER_NOT_FOUND")
	}

	return solved, err
}

func reply(endpoint string, connectionID string, message string) error {
	outgoing := outgoing{
		Message: message,
	}

	data, err := json.Marshal(&outgoing)
	if err != nil {
		return err
	}

	ws.Send(agwSvc, endpoint, []string{connectionID}, data)
	return nil
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
		return response{StatusCode: 500}, errors.New("ROOM_STATUS_INVALID")
	}

	// parse request body
	var incoming incoming
	err = json.Unmarshal([]byte(request.Body), &incoming)
	if err != nil {
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

	// check answer
	isCorrect := checkAnswer(incoming.Problem, incoming.Answer)

	// check challenger status
	challengerSolved, err := checkChallenger(connectionID)
	if err != nil {
		ws.Disconnect(agwSvc, endpoint, connectionID)
		fmt.Println(err)
		return response{StatusCode: 500}, err
	}

	// reply
	if !isCorrect {
		err = reply(endpoint, connectionID, "WRONG_ANSWER")
	} else if challengerSolved {
		err = reply(endpoint, connectionID, "YOU_LOSE")
		ws.Disconnect(agwSvc, endpoint, connectionID)
	} else {
		err = reply(endpoint, connectionID, "YOU_WIN")
		ws.Disconnect(agwSvc, endpoint, connectionID)
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
