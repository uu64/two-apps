package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type Request events.APIGatewayWebsocketProxyRequest
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request Request) (Response, error) {

	var config *aws.Config
	newSession, err := session.NewSession(config)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	svc := apigatewaymanagementapi.New(newSession)
	svc.Endpoint = fmt.Sprintf("https://%s/%s", request.RequestContext.DomainName, request.RequestContext.Stage)
	fmt.Println(svc.Endpoint)

	svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &request.RequestContext.ConnectionID,
		Data:         []byte("Hello"),
	})

	return Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
