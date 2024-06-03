package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	for _, record := range event.Records {
		tales := tales.New[string](tales.Config[string]{
			Source: audio.New(audio.Config{
				Event: record,
				Client: oai.New(oai.Config{
					APIKey:  os.Getenv("OPENAI_API_KEY"),
					BaseURL: "https://api.openai.com",
				}),
			}),
			Destination: destination.Log[string]{},
		})

		err := tales.Run(ctx)
		if err != nil {
			log.Fatalf("Failed to run: %v", err)
		}
	}

	log.Printf("Sent [%d] tales.", len(event.Records))

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
