package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request events.APIGatewayWebsocketProxyRequest
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request Request) (Response, error) {
	fmt.Println("disconnected!!!!!!!")

	return Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
