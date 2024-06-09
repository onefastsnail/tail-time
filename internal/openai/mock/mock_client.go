// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	openai "tail-time/internal/openai"

	gomock "github.com/golang/mock/gomock"
)

// MockClientAPI is a mock of ClientAPI interface.
type MockClientAPI struct {
	ctrl     *gomock.Controller
	recorder *MockClientAPIMockRecorder
}

// MockClientAPIMockRecorder is the mock recorder for MockClientAPI.
type MockClientAPIMockRecorder struct {
	mock *MockClientAPI
}

// NewMockClientAPI creates a new mock instance.
func NewMockClientAPI(ctrl *gomock.Controller) *MockClientAPI {
	mock := &MockClientAPI{ctrl: ctrl}
	mock.recorder = &MockClientAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientAPI) EXPECT() *MockClientAPIMockRecorder {
	return m.recorder
}

// ChatCompletion mocks base method.
func (m *MockClientAPI) ChatCompletion(ctx context.Context, prompt openai.ChatCompletionPrompt) (*openai.ChatCompletionPromptResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatCompletion", ctx, prompt)
	ret0, _ := ret[0].(*openai.ChatCompletionPromptResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChatCompletion indicates an expected call of ChatCompletion.
func (mr *MockClientAPIMockRecorder) ChatCompletion(ctx, prompt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatCompletion", reflect.TypeOf((*MockClientAPI)(nil).ChatCompletion), ctx, prompt)
}

// Completion mocks base method.
func (m *MockClientAPI) Completion(ctx context.Context, prompt openai.CompletionPrompt) (*openai.CompletionPromptResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Completion", ctx, prompt)
	ret0, _ := ret[0].(*openai.CompletionPromptResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Completion indicates an expected call of Completion.
func (mr *MockClientAPIMockRecorder) Completion(ctx, prompt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Completion", reflect.TypeOf((*MockClientAPI)(nil).Completion), ctx, prompt)
}

// TextToSpeech mocks base method.
func (m *MockClientAPI) TextToSpeech(ctx context.Context, prompt openai.TextToSpeechPrompt) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TextToSpeech", ctx, prompt)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TextToSpeech indicates an expected call of TextToSpeech.
func (mr *MockClientAPIMockRecorder) TextToSpeech(ctx, prompt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TextToSpeech", reflect.TypeOf((*MockClientAPI)(nil).TextToSpeech), ctx, prompt)
}
