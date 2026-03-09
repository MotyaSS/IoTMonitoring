package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/processor/kafka/consumer"
	"github.com/MotyaSS/IoTMonitoring/internal/processor/kafka/producer"
	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/MotyaSS/IoTMonitoring/internal/types"
)

type Service struct {
	p   *producer.ProcessorProducer
	c   *consumer.ProcessorConsumer
	log *slog.Logger
}

func NewProcessorService(p *producer.ProcessorProducer, c *consumer.ProcessorConsumer, log *slog.Logger) *Service {
	return &Service{
		p:   p,
		c:   c,
		log: log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := s.c.Consume(ctx)
			if err != nil {
				s.log.Error("failed to consume",
					"err", err,
				)
			}

			var res pb.Telemetry
			err = json.Unmarshal(m.Value, &res)
			if err != nil {
				s.log.Error("failed to unmarshal",
					"err", err,
				)
			}

			data, err := ProcessData(&res)
			if err != nil {
				s.log.Error("failed to process data",
					"err", err,
				)
			}

			err = s.p.Produce(ctx, data)
			if err != nil {
				s.log.Error("failed to produce",
					"err", err,
				)
			}
		}
	}
}

func ProcessData(in *pb.Telemetry) (types.ProcessedData, error) {
	v := types.ProcessedData{
		SenderId:    in.SenderId,
		Timestamp:   in.Timestamp,
		Latitude:    in.Latitude,
		Longitude:   in.Longitude,
		LogMessage:  in.LogMessage,
		Temperature: in.Temperature,
		Pressure:    in.Pressure,
		WindSpeed:   in.WindSpeed,
	}

	// some intense processing
	<-time.After(100 * time.Millisecond)
	return v, nil
}
