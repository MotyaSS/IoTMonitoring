package producer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/segmentio/kafka-go"
)

type ScrapperProducer struct {
	w   *kafka.Writer
	log *slog.Logger
}

func NewScrapperProducer(cfg *config.KafkaConfig, log *slog.Logger) (*ScrapperProducer, error) {
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
	return &ScrapperProducer{
		w:   w,
		log: log,
	}, nil
}

func (kp *ScrapperProducer) Produce(ctx context.Context, msg *pb.Telemetry) error {
	v, err := json.Marshal(msg)
	if err != nil {
		kp.log.Error("Failed to marshal telemetry to json", "error", err)
		return err
	}
	kp.log.Debug("Sending telemetry to kafka", "payload", string(v))
	return kp.w.WriteMessages(
		ctx,
		kafka.Message{
			Key:   nil,
			Value: v,
		},
	)

}

func (kp *ScrapperProducer) Close() error {
	return kp.w.Close()
}
