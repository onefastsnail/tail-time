package audio

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"tail-time/internal/openai"
	"tail-time/internal/source"
	"tail-time/internal/tale"
)

type Config struct {
	OpenAiClient openai.ClientAPI
	TaleClient   source.Source[tale.Tale]
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
	t, err := o.config.TaleClient.Generate(ctx)
	if err != nil {
		return []byte{}, fmt.Errorf("could generate text from source: %v", err)
	}

	voices := []string{"onyx", "echo", "nova"}
	randIndex := rand.Intn(len(voices))

	audio, err := o.convertTextToAudio(ctx, t.Text, voices[randIndex])
	if err != nil {
		return []byte{}, fmt.Errorf("could not convert text to audio: %v", err)
	}

	return audio, nil
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
