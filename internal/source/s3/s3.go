package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Event events.S3EventRecord
}

type S3 struct {
	config Config
}

func New(config Config) *S3 {
	return &S3{config: config}
}

func (s S3) Generate(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "eu-central-1"
	})

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Event.S3.Bucket.Name),
		Key:    aws.String(s.config.Event.S3.Object.Key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object from s3: %w", err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object body from s3: %w", err)
	}

	return string(body), nil
}
