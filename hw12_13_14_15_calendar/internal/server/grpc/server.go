package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
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
}

func NewServer(app *app.App, address string) (*Server, error) {
	app.Log.Info("grpc is running...")
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		app.Log.Error("Listening Error", zap.Error(err))
		return &Server{}, fmt.Errorf("database query failed: %w", err)
	}
	server := grpc.NewServer()
	srv := &Server{ //nolint
		app:    app,
		server: server,
		lsn:    lsn,
	}
	pb.RegisterCalendarServer(server, srv)
	return srv, nil
}

func (s *Server) Start() error {
	s.app.Log.Info("starting grpc server", zap.String("Addr", s.lsn.Addr().String()))
	if err := s.server.Serve(s.lsn); err != nil {
		s.app.Log.Error("Error", zap.Error(err))
		return fmt.Errorf("creating a new ServerTransport failed: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.Id) (*pb.Event, error) {
	var eve storage.Event
	var err error
	s.app.Log.Info("Get Event grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "")
	}
	if !s.app.Mode {
		eve, err = s.app.Get(ctx, req.Id)
		if err != nil {
			s.app.Log.Error("Get Event grpc psql method", zap.Error(err))
		}
	} else {
		eve, err = s.app.GetInMemory(ctx, req.Id)
		if err != nil {
			s.app.Log.Error("Get Event grpc memory method", zap.Error(err))
		}
	}
	startTime, err := ptypes.TimestampProto(eve.StartTime)
	if err != nil {
		s.app.Log.Error("TimestampProto", zap.Error(err))
	}
	endTime, err := ptypes.TimestampProto(eve.EndTime)
	if err != nil {
		s.app.Log.Error("TimestampProto", zap.Error(err))
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
	s.app.Log.Info("Create Event grpc method", zap.Int("req", int(req.Id)))
	if req.Title == "" {
		s.app.Log.Info("BadRequest", zap.Int("Title can't be empty", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "Not enough arguments")
	}
	eve := set(req)

	if !s.app.Mode {
		id, err = s.app.Create(ctx, eve)
		if err != nil {
			s.app.Log.Info("Create Event grpc psql method", zap.String("error", err.Error()))
			return nil, nil
		}
	} else {
		id, err = s.app.CreateInMemory(ctx, eve)
		if err != nil {
			s.app.Log.Info("Create Event grpc memory method", zap.String("error", err.Error()))
			return nil, nil
		}
	}
	return &pb.Id{Id: id}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	var err error
	s.app.Log.Info("Update grpc method", zap.Int("req", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}
	eve := set(req)
	if !s.app.Mode {
		err = s.app.Update(s.ctx, eve)
		if err != nil {
			s.app.Log.Info("Update Event grpc psql method", zap.String("error", err.Error()))
			return nil, nil
		}
	} else {
		err = s.app.UpdateInMemory(s.ctx, eve)
		if err != nil {
			s.app.Log.Info("Update Event grpc memory method", zap.String("error", err.Error()))
			return nil, nil
		}
	}
	return &pb.Id{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.Id) (*emptypb.Empty, error) {
	var err error
	s.app.Log.Info("Delete Event grpc method", zap.Int("id", int(req.Id)))
	if req.Id == 0 {
		s.app.Log.Info("ID can`t be 0", zap.Int("id", int(req.Id)))
		return nil, status.Error(codes.InvalidArgument, "ID can`t be 0")
	}
	if !s.app.Mode {
		err = s.app.Delete(ctx, req.Id)
		if err != nil {
			s.app.Log.Error("Delete Event grpc psql method", zap.Error(err))
		}
	} else {
		err = s.app.DeleteInMemory(ctx, req.Id)
		if err != nil {
			s.app.Log.Error("Delete Event grpc memory method", zap.Error(err))
		}
	}
	return &empty.Empty{}, nil
}

func set(req *pb.Event) storage.Event {
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
