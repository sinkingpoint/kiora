// Code generated by MockGen. DO NOT EDIT.
// Source: ./lib/kiora/kioradb/db.go

// Package mock_kioradb is a generated GoMock package.
package mock_kioradb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/sinkingpoint/kiora/lib/kiora/model"
)

// MockModelReader is a mock of ModelReader interface.
type MockModelReader struct {
	ctrl     *gomock.Controller
	recorder *MockModelReaderMockRecorder
}

// MockModelReaderMockRecorder is the mock recorder for MockModelReader.
type MockModelReaderMockRecorder struct {
	mock *MockModelReader
}

// NewMockModelReader creates a new mock instance.
func NewMockModelReader(ctrl *gomock.Controller) *MockModelReader {
	mock := &MockModelReader{ctrl: ctrl}
	mock.recorder = &MockModelReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelReader) EXPECT() *MockModelReaderMockRecorder {
	return m.recorder
}

// GetAlerts mocks base method.
func (m *MockModelReader) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlerts", ctx)
	ret0, _ := ret[0].([]model.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlerts indicates an expected call of GetAlerts.
func (mr *MockModelReaderMockRecorder) GetAlerts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlerts", reflect.TypeOf((*MockModelReader)(nil).GetAlerts), ctx)
}

// GetExistingAlert mocks base method.
func (m *MockModelReader) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExistingAlert", ctx, labels)
	ret0, _ := ret[0].(*model.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExistingAlert indicates an expected call of GetExistingAlert.
func (mr *MockModelReaderMockRecorder) GetExistingAlert(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExistingAlert", reflect.TypeOf((*MockModelReader)(nil).GetExistingAlert), ctx, labels)
}

// GetSilences mocks base method.
func (m *MockModelReader) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSilences", ctx, labels)
	ret0, _ := ret[0].([]model.Silence)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSilences indicates an expected call of GetSilences.
func (mr *MockModelReaderMockRecorder) GetSilences(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSilences", reflect.TypeOf((*MockModelReader)(nil).GetSilences), ctx, labels)
}

// MockModelWriter is a mock of ModelWriter interface.
type MockModelWriter struct {
	ctrl     *gomock.Controller
	recorder *MockModelWriterMockRecorder
}

// MockModelWriterMockRecorder is the mock recorder for MockModelWriter.
type MockModelWriterMockRecorder struct {
	mock *MockModelWriter
}

// NewMockModelWriter creates a new mock instance.
func NewMockModelWriter(ctrl *gomock.Controller) *MockModelWriter {
	mock := &MockModelWriter{ctrl: ctrl}
	mock.recorder = &MockModelWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelWriter) EXPECT() *MockModelWriterMockRecorder {
	return m.recorder
}

// ProcessAlerts mocks base method.
func (m *MockModelWriter) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range alerts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProcessAlerts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessAlerts indicates an expected call of ProcessAlerts.
func (mr *MockModelWriterMockRecorder) ProcessAlerts(ctx interface{}, alerts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, alerts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessAlerts", reflect.TypeOf((*MockModelWriter)(nil).ProcessAlerts), varargs...)
}

// ProcessSilences mocks base method.
func (m *MockModelWriter) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range silences {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProcessSilences", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessSilences indicates an expected call of ProcessSilences.
func (mr *MockModelWriterMockRecorder) ProcessSilences(ctx interface{}, silences ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, silences...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessSilences", reflect.TypeOf((*MockModelWriter)(nil).ProcessSilences), varargs...)
}

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// GetAlerts mocks base method.
func (m *MockDB) GetAlerts(ctx context.Context) ([]model.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlerts", ctx)
	ret0, _ := ret[0].([]model.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlerts indicates an expected call of GetAlerts.
func (mr *MockDBMockRecorder) GetAlerts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlerts", reflect.TypeOf((*MockDB)(nil).GetAlerts), ctx)
}

// GetExistingAlert mocks base method.
func (m *MockDB) GetExistingAlert(ctx context.Context, labels model.Labels) (*model.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExistingAlert", ctx, labels)
	ret0, _ := ret[0].(*model.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExistingAlert indicates an expected call of GetExistingAlert.
func (mr *MockDBMockRecorder) GetExistingAlert(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExistingAlert", reflect.TypeOf((*MockDB)(nil).GetExistingAlert), ctx, labels)
}

// GetSilences mocks base method.
func (m *MockDB) GetSilences(ctx context.Context, labels model.Labels) ([]model.Silence, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSilences", ctx, labels)
	ret0, _ := ret[0].([]model.Silence)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSilences indicates an expected call of GetSilences.
func (mr *MockDBMockRecorder) GetSilences(ctx, labels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSilences", reflect.TypeOf((*MockDB)(nil).GetSilences), ctx, labels)
}

// ProcessAlerts mocks base method.
func (m *MockDB) ProcessAlerts(ctx context.Context, alerts ...model.Alert) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range alerts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProcessAlerts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessAlerts indicates an expected call of ProcessAlerts.
func (mr *MockDBMockRecorder) ProcessAlerts(ctx interface{}, alerts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, alerts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessAlerts", reflect.TypeOf((*MockDB)(nil).ProcessAlerts), varargs...)
}

// ProcessSilences mocks base method.
func (m *MockDB) ProcessSilences(ctx context.Context, silences ...model.Silence) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range silences {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProcessSilences", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessSilences indicates an expected call of ProcessSilences.
func (mr *MockDBMockRecorder) ProcessSilences(ctx interface{}, silences ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, silences...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessSilences", reflect.TypeOf((*MockDB)(nil).ProcessSilences), varargs...)
}
