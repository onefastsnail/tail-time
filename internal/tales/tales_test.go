package tales

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	mockDestination "tail-time/internal/destination/mock"
	mockSource "tail-time/internal/source/mock"
	"tail-time/internal/tale"
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
	// Given: The source will generate a tale
	t := tale.Tale{Topic: "Test", Text: "Hoho"}

	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Return(t, nil)

	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save(t).Times(1)

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	// When: The tales workload is run
	err := tales.Run(context.TODO())

	// Then: There should be no errors
	s.NoError(err)
}

func (s *TalesSuite) TestSource_Fails() {
	// Given: The source will fail to generate a tale
	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Return(tale.Tale{}, errors.New("oops"))

	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save(gomock.Any()).Times(0)

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	// When: The tales workload is run
	err := tales.Run(context.TODO())

	// Then: An error from the source
	s.Equal("failed to generate tale: oops", err.Error())
}

func (s *TalesSuite) TestDestination_Fails() {
	// Given: The source generates a tale
	source := mockSource.NewMockSource(s.ctrl)
	source.EXPECT().Generate(gomock.Any()).Times(1)

	// And: The destination will fail to save
	destination := mockDestination.NewMockDestination(s.ctrl)
	destination.EXPECT().Save(gomock.Any()).Return(errors.New("oops"))

	tales := New(Config{
		Source:      source,
		Destination: destination,
	})

	// When: The tales workload is run
	err := tales.Run(context.TODO())

	// Then: An error from the destination
	s.Equal("failed to send tale to destination: oops", err.Error())
}
