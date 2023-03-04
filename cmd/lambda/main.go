package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	fmt.Println("event received", event)

	topic := "dinosaurs" // TODO get from event

	tale, err := tales.Generate(ctx, topic)
	if err != nil {
		return "", fmt.Errorf("failed to generate tale: [%w]", err)
	}

	err = tales.SendToKindle(tale)
	if err != nil {
		return "", fmt.Errorf("failed to send tale to Kindle: [%w]", err)
	}

	return fmt.Sprintf("A tale about %s... %s", topic, tale), nil
}

func main() {
	lambda.Start(HandleRequest)
}
