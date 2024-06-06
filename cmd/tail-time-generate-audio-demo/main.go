package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"tail-time/internal/aws"
	"tail-time/internal/destination/localfs"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	record := aws.S3EventDetail{}

	tales := tales.New[[]byte](tales.Config[[]byte]{
		Source: audio.New(audio.Config{
			Event: record,
			Client: oai.New(oai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: localfs.New(localfs.Config{
			Path: "./test.mpga",
		}),
	})

	err = tales.Run(context.TODO())
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
