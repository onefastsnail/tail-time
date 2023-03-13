package destination

import (
	"log"
)

type Log struct{}

func (l Log) Save(tale string) error {
	log.Printf("Got tale... %s", tale)
	return nil
}
