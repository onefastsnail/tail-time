package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"tail-time/internal/destination"
	"tail-time/internal/source/openai"
	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//topic := os.Args[1]
	topic := "dinosaurs"

	tales := tales.New(tales.Config{
		Source:      openai.New(openai.Config{APIKey: os.Getenv("OPENAI_API_KEY")}),
		Destination: destination.Log{},
	})

	tales.Run(context.TODO(), topic)
}
