package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	_aws "tail-time/internal/aws"
	"tail-time/internal/tale"
)

type Config struct {
	Region string
	Event  _aws.S3EventDetail
}

type S3 struct {
	config Config
}

func New(config Config) *S3 {
	return &S3{config: config}
}

func (s S3) Generate(ctx context.Context) (tale.Tale, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = s.config.Region
	})

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Event.Bucket.Name),
		Key:    aws.String(s.config.Event.Object.Key),
	})
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to get object from s3: %w", err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to read object body from s3: %w", err)
	}

	var t tale.Tale
	err = json.Unmarshal(body, &t)
	if err != nil {
		return tale.Tale{}, fmt.Errorf("failed to unmarshal tale: %w", err)
	}

	return t, nil
}
