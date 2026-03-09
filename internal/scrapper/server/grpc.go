package server

import (
	"context"
	"log/slog"
	"net"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	addr string
	log  *slog.Logger
	pb.UnimplementedScrapperServer
}

func NewServer(addr string, log *slog.Logger) *grpcServer {
	if log == nil {
		log = slog.Default()
	}
	return &grpcServer{
		addr: addr,
		log:  log,
	}
}

func (s *grpcServer) SendTelemetry(ctx context.Context, in *pb.Telemetry) (*emptypb.Empty, error) {
	s.log.Info("got telemetry", "id", in.SenderId)
	return nil, nil
}

func (s *grpcServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterScrapperServer(server, s)
	s.log.Info("grpc server listening",
		"address", s.addr,
	)
	return server.Serve(lis)
}
