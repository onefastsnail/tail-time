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

type CompletionPromptResponse struct {
	ID      string `json:"id"`
	Created int    `json:"created"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
}

type client struct {
	apiKey     string
	httpClient http.Client
}

func NewClient(APIKey string) *client {
	return &client{
		apiKey:     APIKey,
		httpClient: http.Client{},
	}
}

func (client *client) Completion(ctx context.Context, prompt CompletionPrompt) (*CompletionPromptResponse, error) {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize prompt: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/completions", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	res, err := client.httpClient.Do(req)
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
