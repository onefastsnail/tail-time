package audio

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	openaimock "tail-time/internal/openai/mock"
	mockSource "tail-time/internal/source/mock"
	"tail-time/internal/tale"
)

type AudioSuite struct {
	suite.Suite
	ctrl *gomock.Controller
}

func TestAudioSuite(t *testing.T) {
	suite.Run(t, new(AudioSuite))
}

func (s *AudioSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *AudioSuite) TestAudio_Generate() {
	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return([]byte("test"), nil)

	mockSourceClient := mockSource.NewMockSource(s.ctrl)
	mockSourceClient.EXPECT().Generate(gomock.Any()).Return(tale.Tale{Text: "test"}, nil)

	audio := New(Config{
		TaleClient:   mockSourceClient,
		OpenAiClient: mockOpenAiClient,
	})

	text, err := audio.Generate(context.Background())

	s.NoError(err)
	s.Equal(text, []byte("test"))
}

func (s *AudioSuite) TestAudio_GenerateTaleFails() {
	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)

	mockSourceClient := mockSource.NewMockSource(s.ctrl)
	mockSourceClient.EXPECT().Generate(gomock.Any()).Return(tale.Tale{}, errors.New("oops"))

	audio := New(Config{
		TaleClient:   mockSourceClient,
		OpenAiClient: mockOpenAiClient,
	})

	text, err := audio.Generate(context.Background())

	s.EqualError(err, "could generate text from source: oops")
	s.Equal(text, []byte{})
}

func (s *AudioSuite) TestAudio_GenerateOpenAIFails() {
	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return(nil, errors.New("oops"))

	mockSourceClient := mockSource.NewMockSource(s.ctrl)
	mockSourceClient.EXPECT().Generate(gomock.Any()).Return(tale.Tale{Text: "test"}, nil)

	audio := New(Config{
		TaleClient:   mockSourceClient,
		OpenAiClient: mockOpenAiClient,
	})

	text, err := audio.Generate(context.Background())

	s.EqualError(err, "could not convert text to audio: failed to get text to speech response: oops")
	s.Equal(text, []byte{})
}
