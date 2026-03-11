package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/MotyaSS/IoTMonitoring/internal/types"
)

type Service struct {
	p   *kafka.Producer
	c   *kafka.Consumer
	log *slog.Logger
}

func NewProcessorService(p *kafka.Producer, c *kafka.Consumer, log *slog.Logger) *Service {
	return &Service{
		p:   p,
		c:   c,
		log: log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.log.Info("starting service")
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
				continue
			}

			var res pb.Telemetry
			err = json.Unmarshal(m.Value, &res)
			if err != nil {
				s.log.Error("failed to unmarshal",
					"err", err,
				)
				continue
			}

			data, err := ProcessData(&res)
			if err != nil {
				s.log.Error("failed to process data",
					"err", err,
				)
				continue
			}

			err = s.p.Produce(ctx, *data)
			if err != nil {
				s.log.Error("failed to produce",
					"err", err,
				)
				continue
			}
		}
	}
}

func ProcessData(in *pb.Telemetry) (*types.ProcessedData, error) {
	// some intense processing
	<-time.After(100 * time.Millisecond)
	t, err := time.Parse(time.DateOnly, in.Timestamp)
	if err != nil {
		return nil, err
	}
	processed := types.ProcessedData{
		SenderId:    in.SenderId,
		Timestamp:   t,
		Latitude:    in.Latitude,
		Longitude:   in.Longitude,
		LogMessage:  in.LogMessage,
		Temperature: in.Temperature,
		Pressure:    in.Pressure,
		WindSpeed:   in.WindSpeed,
	}

	return &processed, nil
}
