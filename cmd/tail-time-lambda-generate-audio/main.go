package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"

	"tail-time/internal/aws"
	s3audio "tail-time/internal/destination/s3/audio"
	"tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/tales"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	var record aws.S3EventDetail
	err := json.Unmarshal(event.Detail, &record)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal event: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %v", err)
	}

	talesWorkload := tales.New[[]byte](tales.Config[[]byte]{
		Source: audio.New(audio.Config{
			Event: record,
			OpenAiClient: openai.New(openai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
			S3ObjectClient: s3.NewFromConfig(cfg, func(s *s3.Options) {
				s.Region = os.Getenv("DESTINATION_BUCKET_REGION")
			}),
		}),
		Destination: s3audio.New(s3audio.Config{
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
