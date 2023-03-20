package tales

import (
	"context"
	"fmt"

	"tail-time/internal/destination"
	"tail-time/internal/source"
)

type Config struct {
	Name        string
	Source      source.Source
	Destination destination.Destination
}

type Tales struct {
	Config
}

func New(config Config) *Tales {
	return &Tales{
		config,
	}
}

func (t Tales) Run(ctx context.Context) error {
	tale, err := t.Source.Generate(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate tale from source: %w", err)
	}

	err = t.Destination.Save(tale)
	if err != nil {
		return fmt.Errorf("failed to send tale to destination: %w", err)
	}

	return nil
}
