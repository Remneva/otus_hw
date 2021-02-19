package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	_ "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

type StoreSuite struct {
	suite.Suite
	mockCtl     *gomock.Controller
	mockDB      *MockEventsStorage
	store       *storage.EventsStorage
	app         *app.App
	start       time.Time
	oneDayLater time.Time
	ctx         context.Context
	log         *zap.Logger
}

func (s *StoreSuite) TeardownTest() {
	s.mockCtl.Finish()
}

func (s *StoreSuite) TestCreate() {
	request := Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	event := storage.Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	jsonBody, _ := json.Marshal(&request)

	handler, _ := newHandler(s.ctx, s.app, s.log)
	ts := httptest.NewServer(handler)
	ts.Close()

	req := httptest.NewRequest("POST", "/set", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()

	s.mockDB.EXPECT().AddEvent(gomock.Any(), event).Return(int64(111), nil)

	handler.setEvent(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	var id JSONID
	err := json.Unmarshal(body, &id)
	if err != nil {
		s.Require().NoError(err)
	}
	var expected JSONID
	expected.ID = 111
	s.Require().Equal(200, resp.Code)
	s.Require().Equal(expected, id)
}

func (s *StoreSuite) TestUpdate() {
	request := Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	event := storage.Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	jsonBody, _ := json.Marshal(&request)
	handler, _ := newHandler(s.ctx, s.app, s.log)
	ts := httptest.NewServer(handler)
	ts.Close()

	req := httptest.NewRequest("POST", "/update", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()

	s.mockDB.EXPECT().UpdateEvent(gomock.Any(), event).Return(nil)

	handler.updateEvent(resp, req)
	s.Require().Equal(200, resp.Code)
}

func (s *StoreSuite) TestUpdateErr() {
	badRequest := Event{
		ID: 0,
	}
	badRequestBody, _ := json.Marshal(&badRequest)
	var errId = errors.New("ID can't be zero or nil value")
	var errorResponse string

	handler, _ := newHandler(s.ctx, s.app, s.log)
	ts := httptest.NewServer(handler)
	ts.Close()

	req := httptest.NewRequest("POST", "/update", bytes.NewBuffer(badRequestBody))
	resp := httptest.NewRecorder()

	handler.updateEvent(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &errorResponse)
	if err != nil {
		s.Require().NoError(err)
	}
	s.Require().Equal(400, resp.Code)
	s.Require().Equal(errId.Error(), errorResponse)
}

func (s *StoreSuite) TestGetEvent() {
	request := int64(111)
	event := storage.Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	jsonBody, _ := json.Marshal(&request)

	handler, _ := newHandler(s.ctx, s.app, s.log)
	ts := httptest.NewServer(handler)
	ts.Close()

	req := httptest.NewRequest("POST", "/get", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()

	s.mockDB.EXPECT().GetEvent(gomock.Any(), int64(111)).Return(event, nil)

	handler.getEvent(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	response := Event{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		s.Require().NoError(err)
	}
	s.Require().Equal(resp.Code, 200)
	s.Require().Equal(event.ID, response.ID)
	s.Require().Equal(event.Title, response.Title)
	s.Require().Equal(event.Description, response.Description)
}

func (s *StoreSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockDB = NewMockEventsStorage(s.mockCtl)
	var z zapcore.Level
	logg, _ := logger.NewLogger(z, "dev", "/dev/null")
	s.log = logg
	s.app = app.NewApp(logg, s.mockDB)
	s.start = time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	s.oneDayLater = s.start.AddDate(0, 0, 1)
}
