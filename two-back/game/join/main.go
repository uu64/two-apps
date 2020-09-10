package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request events.APIGatewayWebsocketProxyRequest
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request Request) (Response, error) {
	fmt.Println("connected!!!!!!!")
	fmt.Println(ctx)
	fmt.Println(request.RequestContext.ConnectionID)

	return createResponse(200, map[string]interface{}{})
}

func createResponse(code int, data map[string]interface{}) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(data)
	if err != nil {
		code = 500
	}
	json.HTMLEscape(&buf, body)

	return Response{
		StatusCode:      code,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, err
}

func main() {
	lambda.Start(Handler)
}
