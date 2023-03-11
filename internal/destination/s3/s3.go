package s3

type Config struct {
	bucket string
	folder string
}

type S3 struct {
	config Config
}

func New(config Config) *S3 {
	return &S3{config: config}
}

func (s S3) Save(tale string) error {
	return nil
}
