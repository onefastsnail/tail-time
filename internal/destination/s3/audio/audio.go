package audio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Config struct {
	Region     string
	BucketName string
	Path       string
}

type S3 struct {
	config Config
}

func New(config Config) *S3 {
	return &S3{config: config}
}

func (s S3) Save(audio []byte) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = s.config.Region
	})

	// TODO: make key same as text counterpart
	objectKey := fmt.Sprintf("%s/%s.mpga", s.config.Path, uuid.New().String())

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(audio),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to s3: %w", err)
	}

	return nil
}
