package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination/s3"
	"tail-time/internal/source/openai"
	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	topic := "dinosaurs" // TODO get from event, Alexa event maybe?

	tales := tales.New(tales.Config{
		Source: openai.New(openai.Config{
			APIKey: os.Getenv("OPENAI_API_KEY"),
			Topic:  topic,
		}),
		Destination: s3.New(s3.Config{
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"), Path: "raw",
		}),
	})

	err := tales.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %e", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
