package destination

import (
	"log"

	"tail-time/internal/tale"
)

type Log struct{}

func (l Log) Save(tale tale.Tale) error {
	log.Printf("Got tale... %+v", tale)
	return nil
}
