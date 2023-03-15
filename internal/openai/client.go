package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CompletionPrompt struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

type CompletionPromptResponseChoice struct {
	Text string `json:"text"`
}

type CompletionPromptResponse struct {
	ID      string                           `json:"id"`
	Created int                              `json:"created"`
	Choices []CompletionPromptResponseChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
}

type Config struct {
	APIKey     string
	HTTPClient http.Client
	BaseURL    string
}

type Client struct {
	config Config
}

func New(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (client *Client) Completion(ctx context.Context, prompt CompletionPrompt) (*CompletionPromptResponse, error) {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize prompt: %w", err)
	}

	url := fmt.Sprintf("%s/v1/completions", client.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.APIKey))
	res, err := client.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	var promptResponse CompletionPromptResponse
	err = json.NewDecoder(res.Body).Decode(&promptResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return &promptResponse, nil
}
