package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"strconv"
	"time"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request Request) (Response, error) {
	q := request.QueryStringParameters
	num := 3
	if v, ok := q["num"]; ok {
		num, _ = strconv.Atoi(v)
	}

	if num > 10 || num < 0 {
		return createResponse(400, map[string]interface{}{
			"message": "invalid parameter!",
		})
	}

	return createResponse(200, map[string]interface{}{
		"problem": createProblem(num),
	})
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

func createProblem(num int) []int {
	rand.Seed(time.Now().UnixNano())

	terms := make([]int, num)
	sum := 2
	for i := 0; i < num - 1; i++ {
		term := rand.Intn(10)
		switch rand.Intn(2) {
		case 0:
			sum = sum + term
		case 1:
			sum = sum - term
		}
		terms[num - 1 - i] = term
	}
	terms[0] = sum

	return terms
}

func main() {
	lambda.Start(Handler)
}
