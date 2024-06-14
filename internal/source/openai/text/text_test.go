package text

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		openAIJSONResponse, _ := json.Marshal(openai.ChatCompletionPromptResponse{Choices: []openai.ChatCompletionPromptChoice{
			{
				Message: openai.ChatCompletionPromptMessage{Content: `{"title": "The Rare Orange Dinosaur", "content": "Once upon a time, there was a test..."}`},
			},
		}})

		_, _ = w.Write(openAIJSONResponse)
	}))
	defer server.Close()

	ai := New(Config{
		Topic:    "dinosaurs",
		Language: "English",
		OpenAiClient: openai.New(openai.Config{
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
