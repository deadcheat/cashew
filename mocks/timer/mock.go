// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/deadcheat/cashew/timer (interfaces: TimeWrapper)

// Package timer is a generated GoMock package.
package timer

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockTimeWrapper is a mock of TimeWrapper interface
type MockTimeWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockTimeWrapperMockRecorder
}

// MockTimeWrapperMockRecorder is the mock recorder for MockTimeWrapper
type MockTimeWrapperMockRecorder struct {
	mock *MockTimeWrapper
}

// NewMockTimeWrapper creates a new mock instance
func NewMockTimeWrapper(ctrl *gomock.Controller) *MockTimeWrapper {
	mock := &MockTimeWrapper{ctrl: ctrl}
	mock.recorder = &MockTimeWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTimeWrapper) EXPECT() *MockTimeWrapperMockRecorder {
	return m.recorder
}

// Now mocks base method
func (m *MockTimeWrapper) Now() time.Time {
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now
func (mr *MockTimeWrapperMockRecorder) Now() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTimeWrapper)(nil).Now))
}
