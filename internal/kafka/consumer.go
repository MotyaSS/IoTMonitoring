package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r   *kafka.Reader
	log *slog.Logger
}

func NewConsumer(cfg *config.KafkaConfig, log *slog.Logger) (*Consumer, error) {
	if cfg.InputTopic == nil {
		return nil, errors.New("input topic required")
	}

	return &Consumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			GroupID: cfg.GroupID,
			Topic:   *cfg.InputTopic,
		}),
		log: log,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) (kafka.Message, error) {
	c.log.Debug("Consuming message at topic" + c.r.Config().Topic)
	return c.r.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
