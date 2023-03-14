package email

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"

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

func writeTaleToDisk(tale string) error {
	file, err := os.Create("/tmp/tale.txt")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(tale)
	if err != nil {
		return fmt.Errorf("failed to write string to file: %w", err)
	}

	return nil
}

func (s Email) Save(tale string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Recipient)
	m.SetHeader("To", s.config.Recipient)
	m.SetHeader("Subject", "A new tale from Tail Time!")
	m.SetBody("text/html", tale)

	err = writeTaleToDisk(tale)
	if err != nil {
		return err
	}

	m.Attach("/tmp/tale.txt")

	var emailRaw bytes.Buffer
	m.WriteTo(&emailRaw)

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
