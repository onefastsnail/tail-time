package tales

import (
	"context"
	"fmt"
	"os"

	"tail-time/internal/openai"
)

//type Tales struct {
//	generator
//}

func Generate(ctx context.Context, topic string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	prompt := openai.CompletionPrompt{
		Model:       "text-davinci-003",
		Prompt:      fmt.Sprintf("Write me a new 100 word story for my kids about %s", topic),
		MaxTokens:   4000,
		Temperature: 0,
	}

	response, err := client.Completion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get prompt completetion: %w", err)
	}

	// TODO handle multiple choices and even none
	return response.Choices[0].Text, nil
}

func SendToKindle(story string) error {
	return nil
}
