package tale

import "time"

type Tale struct {
	Topic     string    `json:"topic"`
	Language  string    `json:"language"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	// TODO Maybe more meta later like model used
}
