package dummy

import (
	"context"
	"fmt"
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

func (d Dummy) Generate(ctx context.Context) (string, error) {
	return fmt.Sprintf("Once upon a time...there was a dummy tale about [%s]. The end.", d.config.Topic), nil
}
