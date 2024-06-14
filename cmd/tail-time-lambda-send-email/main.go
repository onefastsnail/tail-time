package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/aws"
	"tail-time/internal/destination/email"
	"tail-time/internal/source/s3"
	"tail-time/internal/tale"
	"tail-time/internal/worker"
)

func HandleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	var record aws.S3EventDetail
	err := json.Unmarshal(event.Detail, &record)
	if err != nil {
		return "fail", err
	}

	worker := worker.New[tale.Tale](worker.Config[tale.Tale]{
		Source: s3.New(s3.Config{
			Region: os.Getenv("SOURCE_BUCKET_REGION"),
			Event:  record,
		}),
		Destination: email.New(email.Config{
			From: os.Getenv("EMAIL_FROM"),
			To:   os.Getenv("EMAIL_TO"),
		}),
	})

	err = worker.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

	return "OK", nil
}

func main() {
	lambda.Start(HandleRequest)
}
