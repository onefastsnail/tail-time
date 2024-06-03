package destination

import (
	"log"
)

type Log[T any] struct{}

func (l Log[T]) Save(tale T) error {
	log.Printf("Saving... %+v", tale)
	return nil
}
