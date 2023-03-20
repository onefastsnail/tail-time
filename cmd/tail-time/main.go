package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"tail-time/internal/destination/s3"
	oai "tail-time/internal/openai"
	"tail-time/internal/source/openai"
	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//topic := os.Args[1]
	//topic := "dinosaurs and friendship"
	topic := "a dog, an owl and a rabbit"

	tales := tales.New(tales.Config{
		Source: openai.New(openai.Config{
			Topic:    topic,
			Language: "English",
			Client: oai.New(oai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		//Source: dummy.New(dummy.Config{Topic: topic}),
		//Destination: destination.Log{},
		//Destination: email.New(email.Config{Recipient: os.Getenv("EMAIL_DESTINATION")}),
		Destination: s3.New(s3.Config{
			Region:     os.Getenv("DESTINATION_BUCKET_REGION"),
			BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
			Path:       "raw",
		}),
	})

	err = tales.Run(context.TODO())
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
