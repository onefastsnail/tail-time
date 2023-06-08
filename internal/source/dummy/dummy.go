package dummy

import (
	"context"
	"time"

	"tail-time/internal/tale"
)

type Config struct {
	Topic string
}

type Dummy struct {
	config Config
}

func New(config Config) *Dummy {
	return &Dummy{config: config}
}

func (d Dummy) Generate(ctx context.Context) (tale.Tale, error) {
	return tale.Tale{
		Topic:     "dummy",
		Language:  "English",
		Title:     "A story",
		Text:      "Once upon a time... The end.",
		CreatedAt: time.Now(),
	}, nil
}
