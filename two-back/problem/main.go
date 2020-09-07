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

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	q := request.QueryStringParameters
	num := 3
	if v, ok := q["num"]; ok {
		num, _ = strconv.Atoi(v)
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"problem": createProblem(num),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
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
