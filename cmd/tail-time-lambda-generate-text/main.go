package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination/s3/json"
	"tail-time/internal/openai"
	"tail-time/internal/source/openai/text"
	"tail-time/internal/tale"
	"tail-time/internal/worker"
)

type customEvent = map[string]string

// TODO get from an Alexa event

func HandleRequest(ctx context.Context, event customEvent) (string, error) {
	log.Print(event)

	log.Printf("Creating tale about [%s]", event["topic"])

	worker := worker.New[tale.Tale](worker.Config[tale.Tale]{
		Source: text.New(text.Config{
			Topic:    event["topic"],
			Language: "English",
			OpenAiClient: openai.New(openai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: json.New(json.Config{
			Region:     os.Getenv("DESTINATION_BUCKET_REGION"),
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
			Path:       "raw",
		}),
	})

	err := worker.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
