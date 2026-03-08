package server

import (
	"context"
	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type grpcServer struct {
	addr string
	pb.UnimplementedScrapperServer
}

func NewServer(addr string) *grpcServer {
	return &grpcServer{
		addr: addr,
	}
}

func (s *grpcServer) SendTelemetry(ctx context.Context, in *pb.Telemetry) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *grpcServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterScrapperServer(grpcServer, s)
	return grpcServer.Serve(lis)
}
