package openai

import (
	"context"
	"fmt"
	"log"
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
	// Be careful not to use all your tokens per minute quota!

	// First get a title of the story
	titlePrompt := oai.ChatCompletionPrompt{
		Model: "gpt-3.5-turbo",
		Messages: []oai.ChatCompletionPromptMessage{
			{Role: "system", Content: fmt.Sprintf("You are a story writer for young children who writes in %s", o.config.Language)},
			{Role: "user", Content: fmt.Sprintf("Give me a title for an exciting story about %s, without quotes", o.config.Topic)},
		},
		MaxTokens:   500,
		Temperature: 1,
	}

	titleResponse, err := o.config.Client.ChatCompletion(ctx, titlePrompt)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get chat completetion subject: %w", err)
	}

	if len(titleResponse.Choices) < 1 {
		return tale.Tale{}, fmt.Errorf("failed to get a choice from subject response [%+v]", titleResponse)
	}

	title := titleResponse.Choices[0].Message.Content

	log.Printf("Title of the story is: [%s]", title)

	prompt := oai.ChatCompletionPrompt{
		Model: "gpt-3.5-turbo",
		Messages: []oai.ChatCompletionPromptMessage{
			{Role: "user", Content: fmt.Sprintf("Please create a 1000 word engaging bedtime story based on the title provided here: [%s]. The story should be age-appropriate and suitable for children, with a clear beginning, middle, and end. The story should capture the reader's imagination and emotions, with characters that are relatable and memorable. The story's theme or moral should be positive and inspiring, teaching children important lessons about kindness, hope, or perseverance.", title)},
		},
		MaxTokens:   3000,
		Temperature: 1,
	}

	response, err := o.config.Client.ChatCompletion(ctx, prompt)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get prompt completetion: %w", err)
	}

	if len(response.Choices) < 1 {
		return tale.Tale{}, fmt.Errorf("failed to get a choice from response [%+v]", response)
	}

	log.Printf("Recieved a [%d] word story", len(response.Choices[0].Message.Content))

	return tale.Tale{
		Topic:     o.config.Topic,
		Title:     title,
		Text:      response.Choices[0].Message.Content,
		Language:  o.config.Language,
		CreatedAt: time.Now(),
	}, nil
}
