package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"tail-time/internal/tale"
)

type Config struct {
	BucketName string
	Path       string
}

type S3 struct {
	config Config
}

func New(config Config) *S3 {
	return &S3{config: config}
}

func (s S3) Save(tale tale.Tale) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "eu-central-1"
	})

	objectKey := fmt.Sprintf("%s/%s.txt", s.config.Path, uuid.New().String())

	t, err := json.Marshal(tale)
	if err != nil {
		return err
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(t),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to s3: %w", err)
	}

	return nil
}
