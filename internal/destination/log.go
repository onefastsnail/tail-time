package destination

import "fmt"

type Log struct{}

func (l Log) Save(tale string) error {
	fmt.Printf("Got tale... %s", tale)
	return nil
}
