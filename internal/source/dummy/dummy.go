package dummy

import (
	"context"
	"time"

	"tail-time/internal/tale"
)

type Config struct {
	Text string
}

type Dummy struct {
	config Config
}

func New(config Config) *Dummy {
	return &Dummy{config: config}
}

func (d Dummy) Generate(_ context.Context) (tale.Tale, error) {
	return tale.Tale{
		Topic:     "dummy",
		Language:  "English",
		Title:     "A story",
		Text:      d.config.Text,
		CreatedAt: time.Now(),
	}, nil
}
