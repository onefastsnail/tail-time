package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"

	"tail-time/internal/destination"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	record := events.S3EventRecord{
		S3: events.S3Entity{
			Bucket: events.S3Bucket{
				Name: "tales-bucket",
				Arn:  "",
			},
			Object: events.S3Object{
				Key: "/path/to/tale.json",
			},
		},
	}

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

	err = tales.Run(context.TODO())
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
