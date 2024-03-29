// Code generated by MockGen. DO NOT EDIT.
// Source: ./lib/kiora/config/provider.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/sinkingpoint/kiora/lib/kiora/config"
	model "github.com/sinkingpoint/kiora/lib/kiora/model"
)

// MockNotifier is a mock of Notifier interface.
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier.
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance.
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// Name mocks base method.
func (m *MockNotifier) Name() config.NotifierName {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(config.NotifierName)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockNotifierMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockNotifier)(nil).Name))
}

// Notify mocks base method.
func (m *MockNotifier) Notify(ctx context.Context, alerts ...model.Alert) *config.NotificationError {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range alerts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Notify", varargs...)
	ret0, _ := ret[0].(*config.NotificationError)
	return ret0
}

// Notify indicates an expected call of Notify.
func (mr *MockNotifierMockRecorder) Notify(ctx interface{}, alerts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, alerts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockNotifier)(nil).Notify), varargs...)
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// GetNotifiersForAlert mocks base method.
func (m *MockConfig) GetNotifiersForAlert(ctx context.Context, alert *model.Alert) []config.NotifierSettings {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifiersForAlert", ctx, alert)
	ret0, _ := ret[0].([]config.NotifierSettings)
	return ret0
}

// GetNotifiersForAlert indicates an expected call of GetNotifiersForAlert.
func (mr *MockConfigMockRecorder) GetNotifiersForAlert(ctx, alert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifiersForAlert", reflect.TypeOf((*MockConfig)(nil).GetNotifiersForAlert), ctx, alert)
}

// Globals mocks base method.
func (m *MockConfig) Globals() *config.Globals {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Globals")
	ret0, _ := ret[0].(*config.Globals)
	return ret0
}

// Globals indicates an expected call of Globals.
func (mr *MockConfigMockRecorder) Globals() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Globals", reflect.TypeOf((*MockConfig)(nil).Globals))
}

// ValidateData mocks base method.
func (m *MockConfig) ValidateData(ctx context.Context, data config.Fielder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateData", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateData indicates an expected call of ValidateData.
func (mr *MockConfigMockRecorder) ValidateData(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateData", reflect.TypeOf((*MockConfig)(nil).ValidateData), ctx, data)
}

// MockTenanter is a mock of Tenanter interface.
type MockTenanter struct {
	ctrl     *gomock.Controller
	recorder *MockTenanterMockRecorder
}

// MockTenanterMockRecorder is the mock recorder for MockTenanter.
type MockTenanterMockRecorder struct {
	mock *MockTenanter
}

// NewMockTenanter creates a new mock instance.
func NewMockTenanter(ctrl *gomock.Controller) *MockTenanter {
	mock := &MockTenanter{ctrl: ctrl}
	mock.recorder = &MockTenanterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTenanter) EXPECT() *MockTenanterMockRecorder {
	return m.recorder
}

// GetTenant mocks base method.
func (m *MockTenanter) GetTenant(ctx context.Context, data config.Fielder) (config.Tenant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTenant", ctx, data)
	ret0, _ := ret[0].(config.Tenant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTenant indicates an expected call of GetTenant.
func (mr *MockTenanterMockRecorder) GetTenant(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenant", reflect.TypeOf((*MockTenanter)(nil).GetTenant), ctx, data)
}
