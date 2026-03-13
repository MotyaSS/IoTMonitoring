package service

import (
	"context"
	"log/slog"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
)

type producer interface {
	Produce(ctx context.Context, msg any) error
}

type Service struct {
	p   producer
	log *slog.Logger
}

func NewScrapperService(p producer, log *slog.Logger) *Service {
	return &Service{
		p:   p,
		log: log,
	}
}

func (s *Service) SendTelemetry(ctx context.Context, in *pb.Telemetry) error {
	s.log.Debug("Sending telemetry to kafka",
		"sender_id", in.SenderId,
	)
	return s.p.Produce(ctx, in)
}

func (s *Service) Authenticate(token string) error {
	_ = token
	// mock
	return nil
}
