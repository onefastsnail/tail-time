package audio

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"tail-time/internal/aws/s3/mock"
	openaimock "tail-time/internal/openai/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
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

func (s *AudioSuite) TestAudio_ConvertTextToAudio() {
	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return([]byte("test"), nil)

	audio := New(Config{
		OpenAiClient: mockOpenAiClient,
	})

	text, err := audio.convertTextToAudio(context.Background(), "test")

	s.NoError(err)
	s.Equal(text, []byte("test"))
}

func (s *AudioSuite) TestAudio_ConvertTextToAudioFails() {
	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return(nil, errors.New("oops"))

	audio := New(Config{
		OpenAiClient: mockOpenAiClient,
	})

	text, err := audio.convertTextToAudio(context.Background(), "test")

	s.EqualError(err, "failed to get text to speech response: oops")
	s.Equal(text, []byte{})
}

func (s *AudioSuite) TestAudio_GetTextFromS3() {
	getObjectOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(strings.NewReader("{\"text\": \"hello there\"}")),
	}

	mockS3Client := mock.NewMockGetObjectAPI(s.ctrl)
	mockS3Client.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(getObjectOutput, nil)

	audio := New(Config{
		S3ObjectClient: mockS3Client,
	})

	text, err := audio.getTextFromS3(context.Background(), "bucket", "key")

	s.NoError(err)
	s.Equal(text, "hello there")
}

func (s *AudioSuite) TestAudio_GetTextFromS3Fails() {
	mockS3Client := mock.NewMockGetObjectAPI(s.ctrl)
	mockS3Client.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(nil, errors.New("oops"))

	audio := New(Config{
		S3ObjectClient: mockS3Client,
	})

	text, err := audio.getTextFromS3(context.Background(), "bucket", "key")

	s.EqualError(err, "failed to get object from s3: oops")
	s.Equal(text, "")
}

func (s *AudioSuite) TestAudio_Generate() {
	getObjectOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(strings.NewReader("{\"text\": \"hello there\"}")),
	}

	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return([]byte("test"), nil)

	mockS3Client := mock.NewMockGetObjectAPI(s.ctrl)
	mockS3Client.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(getObjectOutput, nil)

	audio := New(Config{
		S3ObjectClient: mockS3Client,
		OpenAiClient:   mockOpenAiClient,
	})

	text, err := audio.Generate(context.Background())

	s.NoError(err)
	s.Equal(text, []byte("test"))
}

func (s *AudioSuite) TestAudio_GenerateGetObjectFails() {
	mockS3Client := mock.NewMockGetObjectAPI(s.ctrl)
	mockS3Client.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(nil, errors.New("oops"))

	audio := New(Config{
		S3ObjectClient: mockS3Client,
	})

	text, err := audio.Generate(context.Background())

	s.EqualError(err, "could not convert text to audio: failed to get object from s3: oops")
	s.Equal(text, []byte{})
}

func (s *AudioSuite) TestAudio_GenerateOpenAIFails() {
	getObjectOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(strings.NewReader("{\"text\": \"hello there\"}")),
	}

	mockOpenAiClient := openaimock.NewMockClientAPI(s.ctrl)
	mockOpenAiClient.EXPECT().TextToSpeech(gomock.Any(), gomock.Any()).Return(nil, errors.New("oops"))

	mockS3Client := mock.NewMockGetObjectAPI(s.ctrl)
	mockS3Client.EXPECT().GetObject(gomock.Any(), gomock.Any()).Return(getObjectOutput, nil)

	audio := New(Config{
		S3ObjectClient: mockS3Client,
		OpenAiClient:   mockOpenAiClient,
	})

	text, err := audio.Generate(context.Background())

	s.EqualError(err, "could not convert text to audio: failed to get text to speech response: oops")
	s.Equal(text, []byte{})
}
