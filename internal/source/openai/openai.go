package openai

import (
	"context"
	"encoding/json"
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

type ChatCompletionPromptStoryResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func New(config Config) *OpenAI {
	return &OpenAI{config: config}
}

func (o OpenAI) Generate(ctx context.Context) (tale.Tale, error) {
	// Be careful not to use all your tokens per minute quota!

	storyPrompt := oai.ChatCompletionPrompt{
		Model: "gpt-3.5-turbo",
		Messages: []oai.ChatCompletionPromptMessage{
			{Role: "system", Content: fmt.Sprintf("You are a story writer for young children who writes in %s", o.config.Language)},
			{Role: "system", Content: "Your stories should be age-appropriate and suitable for children, with a clear beginning, middle, and end. The story should capture the reader's imagination and emotions, with characters that are relatable and memorable. The story's theme or moral should be positive and inspiring, teaching children important lessons about kindness, hope, team work or perseverance."},
			{Role: "system", Content: `You only reply in JSON. The JSON format of your reply should be: {"title": "Story Title", "content": "Story content goes here."}.`},
			{Role: "user", Content: fmt.Sprintf("Write an exciting story about %s", o.config.Topic)},
		},
		MaxTokens:   500,
		Temperature: 1,
	}

	response, err := o.config.Client.ChatCompletion(ctx, storyPrompt)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get chat completetion response: %w", err)
	}

	var promptResponse ChatCompletionPromptStoryResponse
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &promptResponse)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to deserialize chat completetion response: %w", err)
	}

	return tale.Tale{
		Topic:     o.config.Topic,
		Title:     promptResponse.Title,
		Text:      promptResponse.Content,
		Language:  o.config.Language,
		CreatedAt: time.Now(),
	}, nil
}
