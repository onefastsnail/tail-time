package tale

import "time"

type Tale struct {
	Topic     string    `json:"topic"`
	Language  string    `json:"language"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Category  string    `json:"category"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"createdAt"`
	// TODO Maybe more meta later like model used
}
