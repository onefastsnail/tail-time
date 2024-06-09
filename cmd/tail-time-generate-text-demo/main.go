package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"tail-time/internal/destination"
	"tail-time/internal/openai"
	"tail-time/internal/source/openai/text"
	"tail-time/internal/tale"
	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//topic := os.Args[1]
	topic := "bikes and forests"

	tales := tales.New[tale.Tale](tales.Config[tale.Tale]{
		Source: text.New(text.Config{
			Topic:    topic,
			Language: "English",
			Client: openai.New(openai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		//Source: dummy.NewText(dummy.Config{Topic: topic}),
		Destination: destination.Log[tale.Tale]{},
		//Destination: email.New(email.Config{From: os.Getenv("EMAIL_FROM"), To: os.Getenv("EMAIL_TO")}),
		//Destination: s3.NewText(s3.Config{
		//	Region:     os.Getenv("DESTINATION_BUCKET_REGION"),
		//	BucketName: os.Getenv("DESTINATION_BUCKET_NAME"),
		//	Path:       "raw",
		//}),
	})

	err = tales.Run(context.TODO())
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
