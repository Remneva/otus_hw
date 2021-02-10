package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/pb"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	l           *zap.Logger
	start       time.Time
	oneDayLater time.Time
	ctx         context.Context
	srv         Server
	starttime   *timestamppb.Timestamp
	endtime     *timestamppb.Timestamp
}

func (s *StoreSuite) TeardownTest() {
	s.mockCtl.Finish()
}
func (s *StoreSuite) TestCreatEvent() {
	req := &pb.Event{
		Id:          111,
		Title:       "test title",
		Description: "test test test",
		Startdate:   "2020-03-01",
		Starttime:   s.starttime,
		Enddate:     "2020-03-01",
		Endtime:     s.endtime,
	}
	s.mockDB.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Return(req.Id, nil)

	resp, err := s.srv.SetEvent(context.Background(), req)
	if err != nil {
		s.Require().NoError(err)
	}

	s.Require().Equal(resp.Id, req.Id)
}

func (s *StoreSuite) TestGetEvent() {

	event := storage.Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   s.start,
		EndDate:     "2020-03-01",
		EndTime:     s.oneDayLater,
	}
	s.mockDB.EXPECT().GetEvent(gomock.Any(), event.ID).Return(event, nil)

	req := &pb.Id{Id: 111}
	resp, err := s.srv.GetEvent(context.Background(), req)
	if err != nil {
		s.Require().NoError(err)
	}
	s.Require().Equal(resp.Id, event.ID)
	s.Require().Equal(resp.Title, event.Title)
	s.Require().Equal(resp.Description, event.Description)
	s.Require().Equal(resp.Startdate, event.StartDate)
}

func (s *StoreSuite) TestUpdateEvent() {

	req := &pb.Event{
		Id:          111,
		Title:       "test title",
		Description: "test test test",
		Startdate:   "2020-03-01",
		Starttime:   s.starttime,
		Enddate:     "2020-03-01",
		Endtime:     s.endtime,
	}
	s.mockDB.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Return(nil)

	_, err := s.srv.UpdateEvent(context.Background(), req)
	if err != nil {
		s.Require().NoError(err)
	}
	e, _ := status.FromError(err)
	s.Require().Equal(codes.OK, e.Code())
}

func (s *StoreSuite) TestUpdateEventError() {

	req := &pb.Event{
		Id: 0,
	}

	_, err := s.srv.UpdateEvent(context.Background(), req)
	if err != nil {
		s.Require().Error(err)
	}
	e, _ := status.FromError(err)
	s.Require().Equal(codes.InvalidArgument, e.Code())
}

func (s *StoreSuite) TestDeleteEvent() {

	req := &pb.Id{Id: 111}
	s.mockDB.EXPECT().DeleteEvent(gomock.Any(), req.Id).Return(nil)

	_, err := s.srv.DeleteEvent(context.Background(), req)
	if err != nil {
		s.Require().NoError(err)
	}
	e, _ := status.FromError(err)
	s.Require().Equal(codes.OK, e.Code())
}

func (s *StoreSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockDB = NewMockEventsStorage(s.mockCtl)
	var z zapcore.Level
	logg, _ := logger.NewLogger(z, "dev", "/dev/null")
	s.app = app.NewApp(logg, s.mockDB)
	s.srv = Server{app: s.app, log: logg}
	s.start = time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	s.oneDayLater = s.start.AddDate(0, 0, 1)
	s.starttime, _ = ptypes.TimestampProto(s.start)
	s.endtime, _ = ptypes.TimestampProto(s.oneDayLater)
}
