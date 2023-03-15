package openai

import (
	"context"
	"fmt"
	"time"

	oai "tail-time/internal/openai"
	"tail-time/internal/tale"
)

type Config struct {
	Topic    string
	Language string
	Client   *oai.Client
}

type OpenAI struct {
	config Config
}

func New(config Config) *OpenAI {
	return &OpenAI{config: config}
}

func (o OpenAI) Generate(ctx context.Context) (tale.Tale, error) {
	prompt := oai.CompletionPrompt{
		Model:       "text-davinci-003", // TODO change model
		Prompt:      fmt.Sprintf("Write me a brand new exciting 800 word story for my kids about %s in %s", o.config.Topic, o.config.Language),
		MaxTokens:   4000,
		Temperature: 0,
	}

	response, err := o.config.Client.Completion(ctx, prompt)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get prompt completetion: %w", err)
	}

	// TODO handle multiple choices and even none
	return tale.Tale{
		Topic:     o.config.Topic,
		Text:      response.Choices[0].Text,
		Language:  o.config.Language,
		CreatedAt: time.Now(),
	}, nil
}
