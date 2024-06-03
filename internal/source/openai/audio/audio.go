package audio

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	oai "tail-time/internal/openai"
)

type Config struct {
	Event  events.S3EventRecord
	Client *oai.Client
}

type Audio struct {
	config Config
}

type Response struct {
	Path string
}

func New(config Config) *Audio {
	return &Audio{config: config}
}

func (o Audio) Generate(ctx context.Context) (string, error) {
	return "test", nil
}
