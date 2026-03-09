package service

import (
	"context"
	"log/slog"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/kafka/producer"
)

type Service struct {
	p   *producer.ScrapperProducer
	log *slog.Logger
}

func New(p *producer.ScrapperProducer, log *slog.Logger) *Service {
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
