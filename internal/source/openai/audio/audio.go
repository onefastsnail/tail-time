package audio

import (
	"context"
	"fmt"

	"tail-time/internal/aws"
	oai "tail-time/internal/openai"
)

type Config struct {
	Event  aws.S3EventDetail
	Client *oai.Client
}

type Audio struct {
	config Config
}

type Response struct {
	Path string
}

func New(config Config) *Audio {
	return &Audio{config: config}
}

func (o Audio) Generate(ctx context.Context) ([]byte, error) {
	text, err := o.getTextFromS3(ctx, o.config.Event.Bucket.Name, o.config.Event.Object.Key)
	if err != nil {
		return []byte{}, fmt.Errorf("could not convert text to audio: %v", err)
	}

	audio, err := o.convertTextToAudio(ctx, text)
	if err != nil {
		return []byte{}, fmt.Errorf("could not convert text to audio: %v", err)
	}

	return audio, nil
}

func (o Audio) getTextFromS3(ctx context.Context, bucket string, key string) (string, error) {
	text := fmt.Sprintf("Lets get %s from %s", key, bucket)

	return text, nil
}

func (o Audio) convertTextToAudio(ctx context.Context, text string) ([]byte, error) {
	prompt := oai.TextToSpeechPrompt{
		Model: "tts-1",
		Voice: "nova",
		Input: text,
	}

	response, err := o.config.Client.TextToSpeech(ctx, prompt)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get text to speech response: %w", err)
	}

	return response, nil
}
