// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package internalhttp

import (
	context "context"
	storage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	reflect "reflect"
	"time"

	gomock "github.com/golang/mock/gomock"
)

// MockBaseStorage is a mock of BaseStorage interface.
type MockBaseStorage struct {
	ctrl     *gomock.Controller
	recorder *MockBaseStorageMockRecorder
}

// MockBaseStorageMockRecorder is the mock recorder for MockBaseStorage.
type MockBaseStorageMockRecorder struct {
	mock *MockBaseStorage
}

// NewMockBaseStorage creates a new mock instance.
func NewMockBaseStorage(ctrl *gomock.Controller) *MockBaseStorage {
	mock := &MockBaseStorage{ctrl: ctrl}
	mock.recorder = &MockBaseStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBaseStorage) EXPECT() *MockBaseStorageMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockBaseStorage) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", ctx, ev)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockBaseStorageMockRecorder) AddEvent(ctx, ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockBaseStorage)(nil).AddEvent), ctx, ev)
}

// Close mocks base method.
func (m *MockBaseStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBaseStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBaseStorage)(nil).Close))
}

// Connect mocks base method.
func (m *MockBaseStorage) Connect(ctx context.Context, dsn string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", ctx, dsn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockBaseStorageMockRecorder) Connect(ctx, dsn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockBaseStorage)(nil).Connect), ctx, dsn)
}

// DeleteEvent mocks base method.
func (m *MockBaseStorage) DeleteEvent(ctx context.Context, ID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockBaseStorageMockRecorder) DeleteEvent(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockBaseStorage)(nil).DeleteEvent), ctx, ID)
}

// GetEvent mocks base method.
func (m *MockBaseStorage) GetEvent(ctx context.Context, ID int64) (storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", ctx, ID)
	ret0, _ := ret[0].(storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockBaseStorageMockRecorder) GetEvent(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockBaseStorage)(nil).GetEvent), ctx, ID)
}

// GetEvents mocks base method.
func (m *MockBaseStorage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", ctx)
	ret0, _ := ret[0].([]storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockBaseStorageMockRecorder) GetEvents(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockBaseStorage)(nil).GetEvents), ctx)
}

// UpdateEvent mocks base method.
func (m *MockBaseStorage) UpdateEvent(ctx context.Context, ev storage.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, ev)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockBaseStorageMockRecorder) UpdateEvent(ctx, ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockBaseStorage)(nil).UpdateEvent), ctx, ev)
}

// MockEventsStorage is a mock of EventsStorage interface.
type MockEventsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockEventsStorageMockRecorder
}

func (m *MockEventsStorage) ChangeStatusByID(ctx context.Context, id int64) error {
	panic("implement me")
}

func (m *MockEventsStorage) GetEventsByPeriod(ctx context.Context, starttime time.Time, endtime time.Time) ([]storage.Event, error) {
	panic("implement me")
}

func (m *MockEventsStorage) GetStatusByID(ctx context.Context, id int64) (int64, error) {
	panic("implement me")
}

func (m *MockEventsStorage) ChangeStateById(ctx context.Context, id int64) error {
	panic("implement me")
}

func (m *MockEventsStorage) Connect(ctx context.Context, dsn string) error {
	panic("implement me")
}

func (m *MockEventsStorage) Close() error {
	panic("implement me")
}

// MockEventsStorageMockRecorder is the mock recorder for MockEventsStorage.
type MockEventsStorageMockRecorder struct {
	mock *MockEventsStorage
}

// NewMockEventsStorage creates a new mock instance.
func NewMockEventsStorage(ctrl *gomock.Controller) *MockEventsStorage {
	mock := &MockEventsStorage{ctrl: ctrl}
	mock.recorder = &MockEventsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventsStorage) EXPECT() *MockEventsStorageMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockEventsStorage) AddEvent(ctx context.Context, ev storage.Event) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", ctx, ev)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockEventsStorageMockRecorder) AddEvent(ctx, ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventsStorage)(nil).AddEvent), ctx, ev)
}

// DeleteEvent mocks base method.
func (m *MockEventsStorage) DeleteEvent(ctx context.Context, ID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockEventsStorageMockRecorder) DeleteEvent(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockEventsStorage)(nil).DeleteEvent), ctx, ID)
}

// GetEvent mocks base method.
func (m *MockEventsStorage) GetEvent(ctx context.Context, ID int64) (storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", ctx, ID)
	ret0, _ := ret[0].(storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockEventsStorageMockRecorder) GetEvent(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventsStorage)(nil).GetEvent), ctx, ID)
}

// GetEvents mocks base method.
func (m *MockEventsStorage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", ctx)
	ret0, _ := ret[0].([]storage.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockEventsStorageMockRecorder) GetEvents(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventsStorage)(nil).GetEvents), ctx)
}

// UpdateEvent mocks base method.
func (m *MockEventsStorage) UpdateEvent(ctx context.Context, ev storage.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, ev)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventsStorageMockRecorder) UpdateEvent(ctx, ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventsStorage)(nil).UpdateEvent), ctx, ev)
}