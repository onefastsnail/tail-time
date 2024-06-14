package worker

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

type Worker[T any] struct {
	Config[T]
}

func New[T any](config Config[T]) *Worker[T] {
	return &Worker[T]{
		config,
	}
}

func (t Worker[T]) Run(ctx context.Context) error {
	tale, err := t.Source.Generate(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate from source: %w", err)
	}

	err = t.Destination.Save(tale)
	if err != nil {
		return fmt.Errorf("failed to send to destination: %w", err)
	}

	return nil
}
