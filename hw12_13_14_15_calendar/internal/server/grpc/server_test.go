package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//func TestServer(t *testing.T) {
//	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
//	starttime, err := ptypes.TimestampProto(start)
//	oneDayLater := start.AddDate(0, 0, 1)
//	endtime, err := ptypes.TimestampProto(oneDayLater)
//	opts := []grpc.DialOption{
//		grpc.WithInsecure(),
//	}
//	conn, err := grpc.Dial("localhost:50051", opts...)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer conn.Close()
//
//	client := pb.NewCalendarClient(conn)
//	request := &pb.Event{
//		Owner:       1,
//		Title:       "Title",
//		Description: "result",
//		Startdate:   "2020-03-01",
//		Starttime:   starttime,
//		Enddate:     "2020-03-01",
//		Endtime:     endtime,
//	}
//
//	t.Run("Create, update, get, delete event", func(t *testing.T) {
//		respId, err := client.SetEvent(context.Background(), request)
//		if err != nil {
//			fmt.Printf("fail to dial: %v\n", err)
//		}
//		request.Id = respId.Id
//		respId, err = client.UpdateEvent(context.Background(), request)
//		if err != nil {
//			fmt.Printf("fail to dial: %v\n", err)
//		}
//		id := &pb.Id{
//			Id: request.Id,
//		}
//		respEvent, err := client.GetEvent(context.Background(), id)
//		if err != nil {
//			fmt.Printf("fail to dial: %v\n", err)
//		}
//		assert.Equal(t, int64(1), respEvent.Owner)
//		assert.Equal(t, "Title", respEvent.Title)
//		assert.Equal(t, "result", respEvent.Description)
//		assert.Equal(t, "2020-03-01T00:00:00Z", respEvent.Startdate)
//		assert.Equal(t, "2020-03-01T00:00:00Z", respEvent.Enddate)
//
//		respEmpty, err := client.DeleteEvent(context.Background(), id)
//		if err != nil {
//			fmt.Printf("fail")
//		}
//		assert.Equal(t, emptypb.Empty{}, respEmpty)
//
//		request.Id = 0
//		respId, err = client.UpdateEvent(context.Background(), request)
//		e, _ := status.FromError(err)
//		fmt.Println(e.Code().String())
//		assert.Equal(t, codes.InvalidArgument, e.Code())
//
//	})
//
//}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
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

type StoreSuite struct {
	suite.Suite
	mockCtl     *gomock.Controller
	mockDB      *MockEventsStorage
	store       *storage.EventsStorage
	app         *app.App
	start       time.Time
	oneDayLater time.Time
	ctx         context.Context
	srv         Server
	starttime   *timestamppb.Timestamp
	endtime     *timestamppb.Timestamp
}

func (s *StoreSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockDB = NewMockEventsStorage(s.mockCtl)
	var z zapcore.Level
	var c configs.Config
	logg, _ := logger.NewLogger(z, "/dev/null")
	s.app = app.New(logg, s.mockDB, c)
	s.srv = Server{app: s.app}
	s.start = time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	s.oneDayLater = s.start.AddDate(0, 0, 1)
	s.starttime, _ = ptypes.TimestampProto(s.start)
	s.endtime, _ = ptypes.TimestampProto(s.oneDayLater)
}
