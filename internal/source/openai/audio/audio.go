package audio

import (
	"context"
	"fmt"

	"tail-time/internal/aws"
	oai "tail-time/internal/openai"
)

type Config struct {
	Event  aws.S3EventDetail
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
	fmt.Printf("%+v", o.config.Event)

	return "test", nil
}
