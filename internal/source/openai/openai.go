package openai

import (
	"context"
	"fmt"

	oai "tail-time/internal/openai"
)

type Config struct {
	APIKey string
	Topic  string
}

type OpenAI struct {
	config Config
}

func New(config Config) *OpenAI {
	return &OpenAI{config: config}
}

func (o OpenAI) Generate(ctx context.Context) (string, error) {
	client := oai.NewClient(o.config.APIKey)

	prompt := oai.CompletionPrompt{
		Model:       "text-davinci-003", // TODO change model
		Prompt:      fmt.Sprintf("Write me a brand new exciting 800 word story for my kids about %s", o.config.Topic),
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
