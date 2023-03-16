package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"tail-time/internal/openai"
)

type openAISuite struct {
	suite.Suite
}

func TestOpenAISuite(t *testing.T) {
	suite.Run(t, new(openAISuite))
}

func (s *openAISuite) TestGenerate_OK() {
	openAIResponse := openai.CompletionPromptResponse{
		ID:      "1",
		Created: 0,
		Choices: []openai.CompletionPromptResponseChoice{
			{
				Text: "\n\nTitle: The Rare Orange Dinosaur\n\nOnce upon a time, there was a test...",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{},
	}

	openAIJSONResponse, _ := json.Marshal(openAIResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(openAIJSONResponse)
	}))
	defer server.Close()

	ai := New(Config{
		Topic:    "dinosaurs",
		Language: "English",
		Client: openai.New(openai.Config{
			APIKey:  "testing",
			BaseURL: server.URL,
		}),
	})

	actual, err := ai.Generate(context.TODO())
	s.NoError(err)

	s.Equal("dinosaurs", actual.Topic)
	s.Equal("The Rare Orange Dinosaur", actual.Title)
	s.Equal("Once upon a time, there was a test...", actual.Text)
	s.Equal("English", actual.Language)
}
