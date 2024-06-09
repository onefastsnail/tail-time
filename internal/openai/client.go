//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClientAPI interface {
	Completion(ctx context.Context, prompt CompletionPrompt) (*CompletionPromptResponse, error)
	ChatCompletion(ctx context.Context, prompt ChatCompletionPrompt) (*ChatCompletionPromptResponse, error)
	TextToSpeech(ctx context.Context, prompt TextToSpeechPrompt) ([]byte, error)
}

type CompletionPrompt struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

type CompletionPromptResponseChoice struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Index        string `json:"index"`
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

type ChatCompletionPrompt struct {
	Model       string                        `json:"model"`
	Messages    []ChatCompletionPromptMessage `json:"messages"`
	Temperature int                           `json:"temperature"`
	MaxTokens   int                           `json:"max_tokens"`
}

type TextToSpeechPrompt struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

type ChatCompletionPromptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionPromptChoice struct {
	Index        int                         `json:"index"`
	Message      ChatCompletionPromptMessage `json:"message"`
	FinishReason string                      `json:"finish_reason"`
}

type ChatCompletionPromptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionPromptResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int                          `json:"created"`
	Choices []ChatCompletionPromptChoice `json:"choices"`
	Usage   ChatCompletionPromptUsage    `json:"usage"`
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

func (client *Client) ChatCompletion(ctx context.Context, prompt ChatCompletionPrompt) (*ChatCompletionPromptResponse, error) {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize prompt: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", client.config.BaseURL)

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

	var promptResponse ChatCompletionPromptResponse
	err = json.NewDecoder(res.Body).Decode(&promptResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return &promptResponse, nil
}

func (client *Client) TextToSpeech(ctx context.Context, prompt TextToSpeechPrompt) ([]byte, error) {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize prompt: %w", err)
	}

	url := fmt.Sprintf("%s/v1/audio/speech", client.config.BaseURL)

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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read audio data: %v", err)
	}

	return data, nil
}
