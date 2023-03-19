// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/clustering/broadcaster.go

// Package mock_clustering is a generated GoMock package.
package mock_clustering

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/sinkingpoint/kiora/lib/kiora/model"
)

// MockBroadcaster is a mock of Broadcaster interface.
type MockBroadcaster struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcasterMockRecorder
}

// MockBroadcasterMockRecorder is the mock recorder for MockBroadcaster.
type MockBroadcasterMockRecorder struct {
	mock *MockBroadcaster
}

// NewMockBroadcaster creates a new mock instance.
func NewMockBroadcaster(ctrl *gomock.Controller) *MockBroadcaster {
	mock := &MockBroadcaster{ctrl: ctrl}
	mock.recorder = &MockBroadcasterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBroadcaster) EXPECT() *MockBroadcasterMockRecorder {
	return m.recorder
}

// BroadcastAlertAcknowledgement mocks base method.
func (m *MockBroadcaster) BroadcastAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastAlertAcknowledgement", ctx, alertID, ack)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadcastAlertAcknowledgement indicates an expected call of BroadcastAlertAcknowledgement.
func (mr *MockBroadcasterMockRecorder) BroadcastAlertAcknowledgement(ctx, alertID, ack interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastAlertAcknowledgement", reflect.TypeOf((*MockBroadcaster)(nil).BroadcastAlertAcknowledgement), ctx, alertID, ack)
}

// BroadcastAlerts mocks base method.
func (m *MockBroadcaster) BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range alerts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BroadcastAlerts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadcastAlerts indicates an expected call of BroadcastAlerts.
func (mr *MockBroadcasterMockRecorder) BroadcastAlerts(ctx interface{}, alerts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, alerts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastAlerts", reflect.TypeOf((*MockBroadcaster)(nil).BroadcastAlerts), varargs...)
}
