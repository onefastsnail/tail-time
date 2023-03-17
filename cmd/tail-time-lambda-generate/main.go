package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination/s3"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai"
	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	topic := "dinosaurs and cars" // TODO get from event, Alexa event maybe?

	tales := tales.New(tales.Config{
		Source: openai.New(openai.Config{
			Topic:    topic,
			Language: "English",
			Client: oai.New(oai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: s3.New(s3.Config{
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
			Path:       "raw",
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