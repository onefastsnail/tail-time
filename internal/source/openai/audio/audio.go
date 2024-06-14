package audio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	_aws "tail-time/internal/aws"
	_s3 "tail-time/internal/aws/s3"
	"tail-time/internal/openai"
	"tail-time/internal/tale"
)

type Config struct {
	Event          _aws.S3EventDetail
	OpenAiClient   openai.ClientAPI
	S3ObjectClient _s3.GetObjectAPI
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

	voices := []string{"onyx", "echo", "nova"}
	randIndex := rand.Intn(len(voices))

	audio, err := o.convertTextToAudio(ctx, text, voices[randIndex])
	if err != nil {
		return []byte{}, fmt.Errorf("could not convert text to audio: %v", err)
	}

	return audio, nil
}

func (o Audio) getTextFromS3(ctx context.Context, bucket string, key string) (string, error) {
	result, err := o.config.S3ObjectClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object from s3: %w", err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object body from s3: %w", err)
	}

	var t tale.Tale
	err = json.Unmarshal(body, &t)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal tale: %w", err)
	}

	return t.Text, nil
}

func (o Audio) convertTextToAudio(ctx context.Context, text string, voice string) ([]byte, error) {
	prompt := openai.TextToSpeechPrompt{
		Model: "tts-1",
		Voice: voice,
		Input: text,
	}

	log.Printf("Sending [%s] to TTS API", text)

	response, err := o.config.OpenAiClient.TextToSpeech(ctx, prompt)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get text to speech response: %w", err)
	}

	return response, nil
}
