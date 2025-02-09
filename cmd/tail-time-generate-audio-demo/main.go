package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"tail-time/internal/destination/localfs"
	"tail-time/internal/openai"
	"tail-time/internal/source/dummy"
	"tail-time/internal/source/openai/audio"
	"tail-time/internal/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	w := worker.New[[]byte](worker.Config[[]byte]{
		Source: audio.New(audio.Config{
			TaleClient: dummy.New(dummy.Config{Text: "Once upon a time. The end."}),
			OpenAiClient: openai.New(openai.Config{
				APIKey:  os.Getenv("OPENAI_API_KEY"),
				BaseURL: "https://api.openai.com",
			}),
		}),
		Destination: localfs.New(localfs.Config{
			Path: fmt.Sprintf("./tmp-tales/%d.mpga", time.Now().Unix()),
		}),
	})

	err = w.Run(context.TODO())
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
