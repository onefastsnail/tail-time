// Code generated by MockGen. DO NOT EDIT.
// Source: source.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	tale "tail-time/internal/tale"

	gomock "github.com/golang/mock/gomock"
)

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockSource) Generate(ctx context.Context) (tale.Tale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", ctx)
	ret0, _ := ret[0].(tale.Tale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockSourceMockRecorder) Generate(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSource)(nil).Generate), ctx)
}
