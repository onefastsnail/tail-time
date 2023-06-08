package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"tail-time/internal/destination/email"
	"tail-time/internal/source/s3"
	"tail-time/internal/tales"
)

func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	for _, record := range event.Records {
		tales := tales.New(tales.Config{
			Source: s3.New(s3.Config{
				Region: os.Getenv("SOURCE_BUCKET_REGION"),
				Event:  record,
			}),
			Destination: email.New(email.Config{
				From: os.Getenv("EMAIL_FROM"),
				To:   os.Getenv("EMAIL_TO"),
			}),
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
