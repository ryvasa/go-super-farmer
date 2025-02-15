// Code generated by MockGen. DO NOT EDIT.
// Source: utils/glob.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGlobFunc is a mock of GlobFunc interface.
type MockGlobFunc struct {
	ctrl     *gomock.Controller
	recorder *MockGlobFuncMockRecorder
}

// MockGlobFuncMockRecorder is the mock recorder for MockGlobFunc.
type MockGlobFuncMockRecorder struct {
	mock *MockGlobFunc
}

// NewMockGlobFunc creates a new mock instance.
func NewMockGlobFunc(ctrl *gomock.Controller) *MockGlobFunc {
	mock := &MockGlobFunc{ctrl: ctrl}
	mock.recorder = &MockGlobFuncMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobFunc) EXPECT() *MockGlobFuncMockRecorder {
	return m.recorder
}

// Glob mocks base method.
func (m *MockGlobFunc) Glob(pattern string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Glob", pattern)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Glob indicates an expected call of Glob.
func (mr *MockGlobFuncMockRecorder) Glob(pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Glob", reflect.TypeOf((*MockGlobFunc)(nil).Glob), pattern)
}
