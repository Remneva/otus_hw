package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/app"
	srv "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/pb"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedCalendarServer
	app    *app.App
	ctx    context.Context
	server *grpc.Server
	lsn    net.Listener
	log    *zap.Logger
}

var _ srv.Stopper = (*Server)(nil)

func NewServer(app *app.App, l *zap.Logger, address string) (*Server, error) {
	l.Info("grpc is running...")
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		l.Error("Listening Error", zap.Error(err))
		return &Server{}, fmt.Errorf("database query failed: %w", err)
	}
	server := grpc.NewServer()
	srv := &Server{
		app:    app,
		server: server,
		lsn:    lsn,
		log:    l,
	}
	pb.RegisterCalendarServer(server, srv)
	return srv, nil
}

func (s *Server) Start() error {
	s.log.Info("starting grpc server", zap.String("Addr", s.lsn.Addr().String()))
	if err := s.server.Serve(s.lsn); err != nil {
		s.log.Error("Error", zap.Error(err))
		return fmt.Errorf("creating a new ServerTransport failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("grpc server shutdown")
	s.server.GracefulStop()
	return nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.Id) (*pb.Event, error) {
	var eve storage.Event
	var err error
	s.log.Info("Get Event grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "")
	}
	eve, err = s.app.Get(ctx, req.Id)
	if err != nil {
		s.log.Error("Get Event grpc psql method", zap.Error(err))
	}
	startTime, err := ptypes.TimestampProto(eve.StartTime)
	if err != nil {
		s.log.Error("TimestampProto", zap.Error(err))
	}
	endTime, err := ptypes.TimestampProto(eve.EndTime)
	if err != nil {
		s.log.Error("TimestampProto", zap.Error(err))
	}
	return &pb.Event{Id: eve.ID,
		Owner:       eve.Owner,
		Title:       eve.Title,
		Description: eve.Description,
		Startdate:   eve.StartDate,
		Starttime:   startTime,
		Enddate:     eve.EndDate,
		Endtime:     endTime,
	}, nil
}

func (s *Server) SetEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	var id int64
	var err error
	s.log.Info("Create Event grpc method", zap.Int("req", int(req.Id)))
	if req.Title == "" {
		s.log.Info("BadRequest", zap.Int("Title can't be empty", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "Not enough arguments")
	}
	eve := parseToEventStorageStruct(req)

	id, err = s.app.Create(ctx, eve)
	if err != nil {
		s.log.Info("Create Event grpc psql method", zap.String("error", err.Error()))
		return nil, nil
	}
	return &pb.Id{Id: id}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	var err error
	s.log.Info("Update grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "Not enough arguments")
	}
	eve := parseToEventStorageStruct(req)
	err = s.app.Update(s.ctx, eve)
	if err != nil {
		s.log.Info("Update Event grpc psql method", zap.String("error", err.Error()))
		return nil, nil
	}
	return &pb.Id{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.Id) (*emptypb.Empty, error) {
	var err error
	s.log.Info("Delete Event grpc method", zap.Int("id", int(req.Id)))
	if req.Id == 0 {
		s.log.Info("ID can`t be 0", zap.Int("id", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "ID can`t be 0")
	}
	err = s.app.Delete(ctx, req.Id)
	if err != nil {
		s.log.Error("Delete Event grpc psql method", zap.Error(err))
	}
	return &empty.Empty{}, nil
}

func parseToEventStorageStruct(req *pb.Event) storage.Event {
	var eve storage.Event
	eve.ID = req.Id
	eve.Owner = req.Owner
	eve.Title = req.Title
	eve.Description = req.Description
	eve.StartDate = req.Startdate
	eve.EndDate = req.Enddate
	eve.StartTime, _ = ptypes.Timestamp(req.Starttime)
	eve.StartTime, _ = ptypes.Timestamp(req.Endtime)
	return eve
}
