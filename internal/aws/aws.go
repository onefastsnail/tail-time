package aws

type S3EventDetail struct {
	Version string `json:"version"`
	Bucket  struct {
		Name string `json:"name"`
	} `json:"bucket"`
	Object struct {
		Key       string `json:"key"`
		Size      int64  `json:"size"`
		ETag      string `json:"etag"`
		Sequencer string `json:"sequencer"`
	} `json:"object"`
	RequestID string `json:"request-id"`
	Requester string `json:"requester"`
	SourceIP  string `json:"source-ip-address"`
	Reason    string `json:"reason"`
}
