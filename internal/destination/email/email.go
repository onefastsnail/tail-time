package email

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"

	"tail-time/internal/tale"
)

type Config struct {
	From string
	To   string
}

type Email struct {
	config Config
}

var re = regexp.MustCompile("[^a-z0-9]+")

func New(config Config) *Email {
	return &Email{config: config}
}

func (s Email) Save(tale tale.Tale) error {
	log.Printf("Sending story to: [%s]", s.config.To)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load sdk config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	recipients := strings.Split(s.config.To, ",")

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", fmt.Sprintf("%s - From Tail Time! - %s", tale.Title, tale.CreatedAt.Format("01-02-2006")))
	m.SetBody("text/plain", tale.Text)

	// A hack to easily attach the doc to the email, will fix
	fileName, err := writeTaleToDisk(tale)
	if err != nil {
		return err
	}

	m.Attach(fileName)

	var emailRaw bytes.Buffer
	_, err = m.WriteTo(&emailRaw)
	if err != nil {
		return fmt.Errorf("failed to dump message into writer: %w", err)
	}

	input := &ses.SendRawEmailInput{
		Destinations: recipients,
		RawMessage: &types.RawMessage{
			Data: emailRaw.Bytes(),
		},
		Source: aws.String(s.config.From),
	}

	_, err = client.SendRawEmail(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func slugify(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func writeTaleToDisk(tale tale.Tale) (string, error) {
	slugifyTitle := slugify(tale.Title)

	fileName := fmt.Sprintf("/tmp/%s-%s.txt", slugifyTitle, tale.CreatedAt.Format("01-02-2006"))

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
