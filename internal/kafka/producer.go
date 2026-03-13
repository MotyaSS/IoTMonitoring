package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w   *kafka.Writer
	log *slog.Logger
}

func NewProducer(cfg *config.KafkaConfig, log *slog.Logger) (*Producer, error) {
	if cfg.OutputTopic == nil {
		return nil, errors.New("output topic required")
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  *cfg.OutputTopic,
		WriteTimeout:           10 * time.Second,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.RoundRobin{},
	}

	return &Producer{
		w:   w,
		log: log,
	}, nil
}

func (kp *Producer) Produce(ctx context.Context, msg any) error {
	v, err := json.Marshal(msg)
	if err != nil {
		kp.log.Error("producer failed to marshal to json", "error", err)
		return err
	}

	kp.log.Debug("sending telemetry to kafka", "payload", string(v))
	return kp.w.WriteMessages(
		ctx,
		kafka.Message{
			Key:   nil,
			Value: v,
		},
	)
}

func (kp *Producer) Close() error {
	return kp.w.Close()
}
