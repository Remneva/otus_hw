package internalgrpc

import (
	"context"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"net"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedCalendarServer
	app *app.App
	ctx context.Context
}

type Application struct {
}

func New(ctx context.Context, app *app.App) *Service {
	return &Service{
		app: app,
		ctx: ctx,
	}
}

func (s *Service) NewServer(address string) (*Service, error) {
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		s.app.Log.Error("Listening Error", zap.Error(err))
		return &Service{}, errors.Wrap(err, "Database query failed")
	}
	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, s)
	s.app.Log.Info("starting grpc server", zap.String("Addr", lsn.Addr().String()))
	if err := server.Serve(lsn); err != nil {
		s.app.Log.Error("Error", zap.Error(err))
		return &Service{}, errors.Wrap(err, "creating a new ServerTransport failed")
	}
	return &Service{}, nil
}

func (s *Service) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Service) GetEvent(ctx context.Context, req *pb.Id) (*pb.Event, error) {
	s.app.Log.Info("Get Event grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "")
	}
	eve, err := s.app.Repo.GetEvent(ctx, req.Id)
	if err != nil {
		s.app.Log.Error("Get Event grpc method", zap.Error(err))
	}
	StartTime, err := ptypes.TimestampProto(eve.StartTime)
	if err != nil {
		s.app.Log.Error("TimestampProto", zap.Error(err))
	}
	EndTime, err := ptypes.TimestampProto(eve.EndTime)
	if err != nil {
		s.app.Log.Error("TimestampProto", zap.Error(err))
	}
	return &pb.Event{Id: eve.ID,
		Owner:       eve.Owner,
		Title:       eve.Title,
		Description: eve.Description,
		Startdate:   eve.StartDate,
		Starttime:   StartTime,
		Enddate:     eve.EndDate,
		Endtime:     EndTime,
	}, nil
}

func (s *Service) SetEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	s.app.Log.Info("Create Event grpc method", zap.Int("req", int(req.Id)))
	if req.Title == "" {
		s.app.Log.Info("BadRequest", zap.Int("Title can't be empty", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "Not enough arguments")
	}
	var eve storage.Event
	eve.Owner = req.Owner
	eve.Title = req.Title
	eve.Description = req.Description
	eve.StartDate = req.Startdate
	eve.EndDate = req.Enddate
	eve.StartTime, _ = ptypes.Timestamp(req.Starttime)
	eve.StartTime, _ = ptypes.Timestamp(req.Endtime)

	id, err := s.app.Repo.AddEvent(ctx, eve)
	if err != nil {
		s.app.Log.Info("Create Event", zap.String("error", err.Error()))
		return nil, nil
	}
	return &pb.Id{Id: id}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	s.app.Log.Info("Update grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}
	var eve storage.Event
	eve.ID = req.Id
	eve.Owner = req.Owner
	eve.Title = req.Title
	eve.Description = req.Description
	eve.StartDate = req.Startdate
	eve.EndDate = req.Enddate
	eve.StartTime, _ = ptypes.Timestamp(req.Starttime)
	eve.StartTime, _ = ptypes.Timestamp(req.Endtime)

	err := s.app.Repo.UpdateEvent(s.ctx, eve)
	if err != nil {
		s.app.Log.Info("Update Event", zap.String("error", err.Error()))
		return nil, nil
	}
	return &pb.Id{}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.Id) (*emptypb.Empty, error) {
	s.app.Log.Info("Delete Event grpc method", zap.Int("id", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("ID can`t be 0", zap.Int("id", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "ID can`t be 0")
	}
	err := s.app.Repo.DeleteEvent(ctx, req.Id)
	if err != nil {
		s.app.Log.Error("Delete Event grpc method", zap.Error(err))
	}
	return &empty.Empty{}, nil
}
