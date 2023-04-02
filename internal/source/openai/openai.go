package openai

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
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

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func sanitize(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func (o OpenAI) Generate(ctx context.Context) (tale.Tale, error) {
	prompt := oai.CompletionPrompt{
		Model:       "text-davinci-003", // TODO change model
		Prompt:      fmt.Sprintf("Write an exciting 1000 word story for young children about %s in %s. And a title for this story.", o.config.Topic, o.config.Language),
		MaxTokens:   4000,
		Temperature: 1,
	}

	response, err := o.config.Client.Completion(ctx, prompt)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get prompt completetion: %w", err)
	}

	log.Print(response.Choices[0])

	// TODO clean this up, a quick hack to get moving
	splits := strings.SplitN(response.Choices[0].Text, "\n\n", 3)
	title := strings.Replace(splits[1], "Title: ", "", 1)

	return tale.Tale{
		Topic:     o.config.Topic,
		Title:     sanitize(title),
		Text:      splits[2],
		Language:  o.config.Language,
		CreatedAt: time.Now(),
	}, nil
}
