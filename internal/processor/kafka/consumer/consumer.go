package consumer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/segmentio/kafka-go"
)

type ProcessorConsumer struct {
	r   *kafka.Reader
	log *slog.Logger
}

func NewProcessorConsumer(cfg *config.KafkaConfig, log *slog.Logger) (*ProcessorConsumer, error) {
	if cfg.InputTopic == nil {
		return nil, errors.New("input topic required")
	}
	return &ProcessorConsumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			GroupID: cfg.GroupID,
			Topic:   *cfg.InputTopic,
		}),
		log: log,
	}, nil
}

func (p *ProcessorConsumer) Consume(ctx context.Context) (kafka.Message, error) {
	return p.r.ReadMessage(ctx)
}

func (p *ProcessorConsumer) Close() error {
	return p.r.Close()
}
