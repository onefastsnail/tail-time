package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Config struct {
	Recipient string
}

type Email struct {
	config Config
}

func New(config Config) *Email {
	return &Email{config: config}
}

func (s Email) Save(tale string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{s.config.Recipient}, // in sandbox recipient must be verified
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(tale),
				},
			},
			Subject: &types.Content{
				Data: aws.String("A new tale from Tail Time"),
			},
		},
		Source: aws.String(s.config.Recipient),
	}

	_, err = client.SendEmail(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
