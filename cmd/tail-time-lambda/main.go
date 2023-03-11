package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination/s3"
	"tail-time/internal/source/openai"
	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	fmt.Println("event received", event)

	topic := "dinosaurs" // TODO get from event, Alexa event maybe?

	tales := tales.New(tales.Config{
		Source:      openai.New(openai.Config{APIKey: os.Getenv("OPENAI_API_KEY")}),
		Destination: s3.New(s3.Config{}),
	})

	tales.Run(ctx, topic)

	return fmt.Sprintf("A tale about %s... %s", topic, "TODO"), nil
}

func main() {
	lambda.Start(HandleRequest)
}
