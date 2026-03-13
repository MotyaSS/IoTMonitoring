package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var AuthErr = errors.New("scrapper failed to authenticate device")

type Service interface {
	SendTelemetry(context.Context, *pb.Telemetry) error
	Authenticate(token string) error
}

type grpcServer struct {
	addr   string
	log    *slog.Logger
	svc    Service
	server *grpc.Server
	pb.UnimplementedScrapperServer
}

func NewServer(addr string, svc Service, log *slog.Logger) *grpcServer {
	if log == nil {
		log = slog.Default()
	}
	return &grpcServer{
		addr: addr,
		log:  log,
		svc:  svc,
	}
}

func (s *grpcServer) SendTelemetry(ctx context.Context, in *pb.Telemetry) (*emptypb.Empty, error) {
	s.log.Debug("got telemetry", "id", in.SenderId)

	if s.svc == nil {
		s.log.Error("no service configured to handle telemetry")
		return nil, nil
	}

	if s.svc.Authenticate(in.AuthToken) != nil {
		s.log.Error("failed to authenticate",
			"sender_id", in.SenderId,
			"token", in.AuthToken,
		)
		return nil, fmt.Errorf("%w: auth token - %s", AuthErr, in.AuthToken)
	}

	if err := s.svc.SendTelemetry(ctx, in); err != nil {
		s.log.Error("failed to forward telemetry to service", "err", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *grpcServer) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	pb.RegisterScrapperServer(s.server, s)
	s.log.Info("grpc server listening",
		"address", s.addr,
	)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.server.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		s.log.Info("shutdown requested, gracefully stopping grpc server")
		s.server.GracefulStop()
		return nil
	case err := <-serveErr:
		return err
	}
}
