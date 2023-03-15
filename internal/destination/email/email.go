package email

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"

	"tail-time/internal/tale"
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

func writeTaleToDisk(tale tale.Tale) (string, error) {
	fileName := fmt.Sprintf("/tmp/a-tale-about-%s-in-%s-%s.txt", tale.Topic, tale.Language, tale.CreatedAt.Format("01-02-2006"))

	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(tale.Text)
	if err != nil {
		return "", fmt.Errorf("failed to write string to file: %w", err)
	}

	return fileName, nil
}

func (s Email) Save(tale tale.Tale) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Recipient)
	m.SetHeader("To", s.config.Recipient)
	m.SetHeader("Subject", fmt.Sprintf("%s - A new tale about %s in %s from Tail Time!", tale.CreatedAt.Format("01-02-2006"), tale.Topic, tale.Language))
	m.SetBody("text/plain", tale.Text)

	fileName, err := writeTaleToDisk(tale)
	if err != nil {
		return err
	}

	m.Attach(fileName)

	var emailRaw bytes.Buffer
	_, err = m.WriteTo(&emailRaw)
	if err != nil {
		return err
	}

	input := &ses.SendRawEmailInput{
		Destinations: []string{s.config.Recipient},
		RawMessage: &types.RawMessage{
			Data: emailRaw.Bytes(),
		},
		Source: aws.String(s.config.Recipient),
	}

	_, err = client.SendRawEmail(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
