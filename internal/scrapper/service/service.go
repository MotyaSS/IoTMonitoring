package service

import (
	"context"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/kafka/producer"
)

type Service struct {
	p *producer.KafkaProducer
}

func New(p *producer.KafkaProducer) *Service {
	return &Service{
		p: p,
	}
}

func (s *Service) SendTelemetry(ctx context.Context, in *pb.Telemetry) error {
	if s == nil || s.p == nil {
		return nil
	}
	return s.p.SendTelemetry(ctx, in)
}
