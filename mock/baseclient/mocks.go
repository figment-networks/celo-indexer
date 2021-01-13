// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/figment-networks/celo-indexer/client (interfaces: Client,RequestCounter)

// Package mock_client is a generated GoMock package.
package mock_client

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// GetName mocks base method
func (m *MockClient) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockClientMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockClient)(nil).GetName))
}

// MockRequestCounter is a mock of RequestCounter interface
type MockRequestCounter struct {
	ctrl     *gomock.Controller
	recorder *MockRequestCounterMockRecorder
}

// MockRequestCounterMockRecorder is the mock recorder for MockRequestCounter
type MockRequestCounterMockRecorder struct {
	mock *MockRequestCounter
}

// NewMockRequestCounter creates a new mock instance
func NewMockRequestCounter(ctrl *gomock.Controller) *MockRequestCounter {
	mock := &MockRequestCounter{ctrl: ctrl}
	mock.recorder = &MockRequestCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestCounter) EXPECT() *MockRequestCounterMockRecorder {
	return m.recorder
}

// GetCounter mocks base method
func (m *MockRequestCounter) GetCounter() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetCounter indicates an expected call of GetCounter
func (mr *MockRequestCounterMockRecorder) GetCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockRequestCounter)(nil).GetCounter))
}

// IncrementCounter mocks base method
func (m *MockRequestCounter) IncrementCounter() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementCounter")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// IncrementCounter indicates an expected call of IncrementCounter
func (mr *MockRequestCounterMockRecorder) IncrementCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementCounter", reflect.TypeOf((*MockRequestCounter)(nil).IncrementCounter))
}

// InitCounter mocks base method
func (m *MockRequestCounter) InitCounter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitCounter")
}

// InitCounter indicates an expected call of InitCounter
func (mr *MockRequestCounterMockRecorder) InitCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitCounter", reflect.TypeOf((*MockRequestCounter)(nil).InitCounter))
}