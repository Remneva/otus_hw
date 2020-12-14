package internalgrpc

import (
	"context"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Service struct {
	pb.UnimplementedCalendarServer
	grpc *grpc.Server
	l    *zap.Logger
}

func NewServer(logger *zap.Logger) (*grpc.Server, error) {
	lsn, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, new(Service))

	logger.Info("starting grpc server", zap.String("Addr", lsn.Addr().String()))
	if err := server.Serve(lsn); err != nil {
		log.Println(err)
		return nil, err
	}
	return server, nil
}

func (s *Service) GetEvent(ctx context.Context, req *pb.Id) (*pb.Event, error) {
	if req.Id == 0 {
		log.Printf("Not enough arguments")
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}
	return &pb.Event{}, nil
}

func (s *Service) SetEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	log.Printf("new event receive (requset=%s, time=%v)",
		req.String(), ptypes.TimestampString(req.Startdate))

	if req.Title == "" {
		log.Printf("Not enough arguments")
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}

	return &pb.Id{}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.Id, error) {
	log.Printf("new event update (requset=%s, time=%v)",
		req.String(), ptypes.TimestampString(req.Startdate))

	if req.Title == "" {
		log.Printf("Not enough arguments")
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}

	log.Printf("vote accepted")
	return &pb.Id{}, nil
}
