package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"tail-time/internal/aws"
	audiod "tail-time/internal/destination/s3/audio"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/tales"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	var record aws.S3EventDetail
	err := json.Unmarshal(event.Detail, &record)
	if err != nil {
		return "fail", fmt.Errorf("failed to unmarshal event: %v", err)
	}

	talesWorkload := tales.New[[]byte](tales.Config[[]byte]{
		Source: audio.New(audio.Config{
			Event: record,
			Client: oai.New(oai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: audiod.New(audiod.Config{
			Region:     os.Getenv("DESTINATION_BUCKET_REGION"),
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
			Path:       "audio",
		}),
	})

	err = talesWorkload.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
