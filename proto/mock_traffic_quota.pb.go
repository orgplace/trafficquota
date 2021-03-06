// Code generated by MockGen. DO NOT EDIT.
// Source: traffic_quota.pb.go

package proto

import (
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockTrafficQuotaClient is a mock of TrafficQuotaClient interface
type MockTrafficQuotaClient struct {
	ctrl     *gomock.Controller
	recorder *MockTrafficQuotaClientMockRecorder
}

// MockTrafficQuotaClientMockRecorder is the mock recorder for MockTrafficQuotaClient
type MockTrafficQuotaClientMockRecorder struct {
	mock *MockTrafficQuotaClient
}

// NewMockTrafficQuotaClient creates a new mock instance
func NewMockTrafficQuotaClient(ctrl *gomock.Controller) *MockTrafficQuotaClient {
	mock := &MockTrafficQuotaClient{ctrl: ctrl}
	mock.recorder = &MockTrafficQuotaClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrafficQuotaClient) EXPECT() *MockTrafficQuotaClientMockRecorder {
	return m.recorder
}

// Take mocks base method
func (m *MockTrafficQuotaClient) Take(ctx context.Context, in *TakeRequest, opts ...grpc.CallOption) (*TakeResponse, error) {
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Take", varargs...)
	ret0, _ := ret[0].(*TakeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Take indicates an expected call of Take
func (mr *MockTrafficQuotaClientMockRecorder) Take(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Take", reflect.TypeOf((*MockTrafficQuotaClient)(nil).Take), varargs...)
}

// MockTrafficQuotaServer is a mock of TrafficQuotaServer interface
type MockTrafficQuotaServer struct {
	ctrl     *gomock.Controller
	recorder *MockTrafficQuotaServerMockRecorder
}

// MockTrafficQuotaServerMockRecorder is the mock recorder for MockTrafficQuotaServer
type MockTrafficQuotaServerMockRecorder struct {
	mock *MockTrafficQuotaServer
}

// NewMockTrafficQuotaServer creates a new mock instance
func NewMockTrafficQuotaServer(ctrl *gomock.Controller) *MockTrafficQuotaServer {
	mock := &MockTrafficQuotaServer{ctrl: ctrl}
	mock.recorder = &MockTrafficQuotaServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrafficQuotaServer) EXPECT() *MockTrafficQuotaServerMockRecorder {
	return m.recorder
}

// Take mocks base method
func (m *MockTrafficQuotaServer) Take(arg0 context.Context, arg1 *TakeRequest) (*TakeResponse, error) {
	ret := m.ctrl.Call(m, "Take", arg0, arg1)
	ret0, _ := ret[0].(*TakeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Take indicates an expected call of Take
func (mr *MockTrafficQuotaServerMockRecorder) Take(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Take", reflect.TypeOf((*MockTrafficQuotaServer)(nil).Take), arg0, arg1)
}
