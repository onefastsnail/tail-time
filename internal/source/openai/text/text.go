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

	length := 1000

	storyPrompt := openai.ChatCompletionPrompt{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionPromptMessage{
			{Role: "system", Content: "You are a skilled children's story writer, crafting engaging, age-appropriate bedtime stories that spark imagination and emotions."},
			{Role: "system", Content: "Each story should have a clear beginning, middle, and end, with memorable characters and a heartwarming or adventurous plot. Include some or all of the following positive morals: kindness, perseverance, teamwork, and hope."},
			{Role: "system", Content: "Use rich sensory descriptions, simple yet engaging language, and dialogue to bring the story to life. The tone should be warm and immersive."},
			{Role: "system", Content: "Provide a category and a one-line summary. Additionally, include a short moral lesson at the end."},
			{Role: "system", Content: `You must reply in valid JSON format. The JSON format of your reply should be: {"title": "Story Title", "content": "Story content goes here.", "category": "Category of the story", "summary": "Summary of the story"}.`},
			{Role: "user", Content: fmt.Sprintf("Write a bedtime story in %s about %s. The story should be approximately %d words long.", o.config.Language, o.config.Topic, length)},
		},
		MaxCompletionTokens: length * 2,
		Temperature:         0.8,
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
