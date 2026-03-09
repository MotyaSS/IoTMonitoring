package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	IoTMonitoring "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	if cfg.OutputTopic == nil {

	}
	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  *cfg.OutputTopic,
		WriteTimeout:           10 * time.Second,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.RoundRobin{},
	}
	return &KafkaProducer{
		w: w,
	}, nil
}

func (kp *KafkaProducer) SendTelemetry(ctx context.Context, msg *IoTMonitoring.Telemetry) error {
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return kp.w.WriteMessages(
		ctx,
		kafka.Message{
			Key:   nil,
			Value: v,
		},
	)

}

func (kp *KafkaProducer) Close() error {
	return kp.w.Close()
}
