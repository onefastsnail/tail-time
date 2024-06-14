package text

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tail-time/internal/openai"
	"tail-time/internal/tale"
)

type Config struct {
	Topic        string
	Language     string
	OpenAiClient openai.ClientAPI
}

type Text struct {
	config Config
}

type ChatCompletionPromptStoryResponse struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

func New(config Config) *Text {
	return &Text{config: config}
}

func (o Text) Generate(ctx context.Context) (tale.Tale, error) {
	// Be careful not to use all your tokens per minute quota!

	storyPrompt := openai.ChatCompletionPrompt{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionPromptMessage{
			{Role: "system", Content: "You are a story writer for young children"},
			{Role: "system", Content: "Your stories should be age-appropriate and suitable for children, with a clear beginning, middle, and end. The story should capture the reader's imagination and emotions, with characters that are relatable and memorable. The story's theme or moral should be positive and inspiring, teaching children important lessons about kindness, hope, team work or perseverance."},
			{Role: "system", Content: "Also provide a category and a one line summary about your stories."},
			{Role: "system", Content: `You only reply in JSON. The JSON format of your reply should be: {"title": "Story Title", "content": "Story content goes here.", "category": "Category of the story", "summary": "Summary of the story"}.`},
			{Role: "user", Content: fmt.Sprintf("Write a 1000 word story in %s about %s", o.config.Language, o.config.Topic)},
		},
		MaxTokens:   1200,
		Temperature: 1,
	}

	response, err := o.config.OpenAiClient.ChatCompletion(ctx, storyPrompt)
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
		Category:  promptResponse.Category,
		Summary:   promptResponse.Summary,
		Language:  o.config.Language,
		CreatedAt: time.Now(),
	}, nil
}
