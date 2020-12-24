package memorystorage

import (
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

//func TestSetEvent(t *testing.T) {
//	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
//	oneDayLater := start.AddDate(0, 0, 1)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	mockCtl := gomock.NewController(t)
//	defer mockCtl.Finish()
//
//	mockDB := internalhttp.NewMockBaseStorage(mockCtl)
//	store := New(mockDB)
//
//	event := sqlstorage.Event{
//		ID:          111,
//		Title:       "test title",
//		Description: "test test test",
//		StartDate:   "2020-03-01",
//		StartTime:   start,
//		EndDate:     "2020-03-01",
//		EndTime:     oneDayLater,
//	}
//
//	mockDB.EXPECT().Insert(ctx, eventMatcher{event}).Return(nil)
//	mockDB.EXPECT().GetLastId(ctx).Return(event.ID, nil)
//
//	newID, err := store.AddEvent(ctx, event)
//
//	require.NoError(t, err)
//	require.NotEqual(t, event.ID, newID)
//}

//func TestStorage(t *testing.T) {
//
//	var z zapcore.Level
//	flag.Parse()
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	storage := new(Storage)
//	var ev sqlstorage.Event
//	logg, err := logger.NewLogger(z, "/dev/null")
//	if err != nil {
//		log.Fatal("failed to create logger")
//	}
//	err = storage.Connect(ctx, dsn, logg)
//	if err != nil {
//		log.Fatal("failed to connect db")
//	}
//
//	t.Run("CRUD query", func(t *testing.T) {
//		id, err := storage.AddEvent(ctx, ev)
//		require.Errorf(t, err, "Database query failed")
//		require.Equal(t, id, int64(0))
//
//		err = storage.UpdateEvent(ctx, ev)
//		require.Errorf(t, err, "Database query failed")
//
//		err = storage.DeleteEvent(ctx, 0)
//		require.NoError(t, err)
//	})
//
//}

//
//type StoreSuite struct {
//	suite.Suite
//	mockCtl *gomock.Controller
//	mockDB  *MockEventDB
//	store   *sqlstorage.EventsStorage
//}
//
//type MockEventDB struct {
//	ctrl     *gomock.Controller
//	recorder *MockEventDBMockRecorder
//}
//
//type MockEventDBMockRecorder struct {
//	mock *MockEventDB
//}
//
//// NewMockUsersDB creates a new mock instance
//func NewMockEventDB(ctrl *gomock.Controller) *MockEventDB {
//	mock := &MockEventDB{ctrl: ctrl}
//	mock.recorder = &MockEventDBMockRecorder{mock}
//	return mock
//}
//
//// EXPECT returns an object that allows the caller to indicate expected use
//func (m *MockEventDB) EXPECT() *MockEventDBMockRecorder {
//	return m.recorder
//}
//
type eventMatcher struct {
	sqlstorage.Event
}

//
//func TestGetEvent(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	mockCtl := gomock.NewController(t)
//	defer mockCtl.Finish()
//
//	mockDB := NewMockEventDB(mockCtl)
//	store := New(mockDB)
//
//	event := sqlstorage.Event{
//		ID:          111,
//		Title:       "title_name",
//		Description: "test test test",
//	}
//
//	mockDB.EXPECT().AddEvent(eventMatcher{event}).Return(event.ID, nil)
//	mockDB.EXPECT().GetEvent(event.ID).Return(event, nil)
//
//	newID, err := store.GetEvent(ctx, event.ID)
//
//	require.NoError(t, err)
//	require.NotEqual(t, newID, event.ID)
//}
//
//// AddEvent mocks base method
//func (m *MockEventDB) AddEvent(ctx context.Context, arg0 sqlstorage.Event) (int64, error) {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "AddEvent", arg0)
//	ret0, _ := ret[0].(int64)
//	ret1, _ := ret[1].(error)
//	return ret0, ret1
//}
//
//// AddEvent indicates an expected call of AddUser
//func (mr *MockEventDBMockRecorder) AddEvent(arg0 interface{}) *gomock.Call {
//	mr.mock.ctrl.T.Helper()
//	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventDB)(nil).AddEvent), arg0)
//}
//
//// AddEvent mocks base method
//func (m *MockEventDB) DeleteEvent(ctx context.Context, arg0 int64) error {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "DeleteEvent", arg0)
//	ret0, _ := ret[0].(error)
//	return ret0
//}
//
//// AddEvent indicates an expected call of AddUser
//func (mr *MockEventDBMockRecorder) DeleteEvent(arg0 interface{}) *gomock.Call {
//	mr.mock.ctrl.T.Helper()
//	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventDB)(nil).DeleteEvent), arg0)
//}
//
//// AddEvent mocks base method
//func (m *MockEventDB) GetEvent(ctx context.Context, arg0 int64) (sqlstorage.Event, error) {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "GetEvent", arg0)
//	ret0, _ := ret[0].(sqlstorage.Event)
//	ret1, _ := ret[1].(error)
//	return ret0, ret1
//}
//
//// AddEvent indicates an expected call of AddUser
//func (mr *MockEventDBMockRecorder) GetEvent(arg0 interface{}) *gomock.Call {
//	mr.mock.ctrl.T.Helper()
//	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventDB)(nil).GetEvent), arg0)
//}
//
//// AddEvent mocks base method
//func (m *MockEventDB) UpdateEvent(ctx context.Context, arg0 sqlstorage.Event) error {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "UpdateEvent", arg0)
//	ret0, _ := ret[0].(error)
//	return ret0
//}
//
//// AddEvent indicates an expected call of AddUser
//func (mr *MockEventDBMockRecorder) UpdateEvent(arg0 interface{}) *gomock.Call {
//	mr.mock.ctrl.T.Helper()
//	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventDB)(nil).UpdateEvent), arg0)
//}
//
//// AddEvent mocks base method
//func (m *MockEventDB) GetEvents(ctx context.Context) ([]sqlstorage.Event, error) {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "GetEvents")
//	ret0, _ := ret[0].([]sqlstorage.Event)
//	ret1, _ := ret[1].(error)
//	return ret0, ret1
//}
//
//// AddEvent indicates an expected call of AddUser
//func (mr *MockEventDBMockRecorder) GetEvents(arg0 interface{}) *gomock.Call {
//	mr.mock.ctrl.T.Helper()
//	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventDB)(nil).GetEvents), arg0)
//}
