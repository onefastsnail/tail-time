package tales

import (
	"context"
	"fmt"

	"tail-time/internal/destination"
	"tail-time/internal/source"
)

type Config[T any] struct {
	Name        string
	Source      source.Source[T]
	Destination destination.Destination[T]
}

type Tales[T any] struct {
	Config[T]
}

func New[T any](config Config[T]) *Tales[T] {
	return &Tales[T]{
		config,
	}
}

func (t Tales[T]) Run(ctx context.Context) error {
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
