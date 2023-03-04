package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"

	"tail-time/internal/tales"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//topic := os.Args[1]
	topic := "dinosaurs"

	tale, err := tales.Generate(context.TODO(), topic)
	if err != nil {
		log.Fatal("Failed to generate tale: [%w]", err)
	}

	err = tales.SendToKindle(tale)
	if err != nil {
		log.Fatalf("Failed to send tale to Kindle: [%e]", err)
	}

	log.Printf("A tale about %s... %s", topic, tale)
}
