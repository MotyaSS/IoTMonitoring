package service

import (
	"context"
	"log/slog"

	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
)

type Service struct {
	p   *kafka.Producer
	log *slog.Logger
}

func NewScrapperService(p *kafka.Producer, log *slog.Logger) *Service {
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
	// mock
	return nil
}
