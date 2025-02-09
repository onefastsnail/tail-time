package localtale

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"tail-time/internal/tale"
)

type Config struct {
	Path string
}

type LocalTale struct {
	config Config
}

func New(config Config) *LocalTale {
	return &LocalTale{
		config: config,
	}
}

func (l LocalTale) Save(tale tale.Tale) error {
	data, err := json.Marshal(tale)

	err = os.MkdirAll(l.config.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("%s/%d.txt", l.config.Path, time.Now().Unix())

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	log.Printf("Saved data [%s] to local file system", filename)

	return nil
}
