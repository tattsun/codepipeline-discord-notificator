package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"string"`
}

func HandleRequest(ctx context.Context, event Event) (Response, error) {
	return Response{
		Message: fmt.Sprintf("Hello, %+v", event),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
