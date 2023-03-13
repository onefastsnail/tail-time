package tales

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	mockDestination "tail-time/internal/destination/mock"
	mockSource "tail-time/internal/source/mock"
)

type TalesSuite struct {
	suite.Suite
	ctrl *gomock.Controller
}

func TestTalesSuite(t *testing.T) {
	suite.Run(t, new(TalesSuite))
}

func (s *TalesSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *TalesSuite) TestRun_OK() {
	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Return("a tale", nil)

	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save("a tale").Times(1)

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	err := tales.Run(context.TODO())

	s.NoError(err)
}

func (s *TalesSuite) TestSource_Fails() {
	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Return("", errors.New("oops"))

	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save(gomock.Any()).Times(0)

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	err := tales.Run(context.TODO())

	s.Equal("failed to generate tale: oops", err.Error())
}

func (s *TalesSuite) TestDestination_Fails() {
	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Times(1)

	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save(gomock.Any()).Return(errors.New("oops"))

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	err := tales.Run(context.TODO())

	s.Equal("failed to send tale to destination: oops", err.Error())
}
