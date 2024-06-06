package localfs

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Path string
}

type LocalFs struct {
	config Config
}

func New(config Config) *LocalFs {
	return &LocalFs{
		config: config,
	}
}

func (l LocalFs) Save(data []byte) error {
	err := os.WriteFile(l.config.Path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	log.Printf("Saved data [%s] to local file system", l.config.Path)

	return nil
}
