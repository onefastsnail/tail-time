package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/aws"
	s3audio "tail-time/internal/destination/s3/audio"
	"tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	s3source "tail-time/internal/source/s3"
	"tail-time/internal/worker"
)

func HandleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	var record aws.S3EventDetail
	err := json.Unmarshal(event.Detail, &record)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal event: %v", err)
	}

	w := worker.New[[]byte](worker.Config[[]byte]{
		Source: audio.New(audio.Config{
			TaleClient: s3source.New(s3source.Config{
				Region: os.Getenv("SOURCE_BUCKET_REGION"),
				Event:  record,
			}),
			OpenAiClient: openai.New(openai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: s3audio.New(s3audio.Config{
			Region:     os.Getenv("DESTINATION_BUCKET_REGION"),
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
			Path:       "audio",
		}),
	})

	err = w.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
