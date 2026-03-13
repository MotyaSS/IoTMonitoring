package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	"github.com/MotyaSS/IoTMonitoring/internal/types"
)

type DeviceStore interface {
	GetDevice(ctx context.Context, deviceID string) (*types.Device, error)
}

type Service struct {
	p          *kafka.Producer
	c          *kafka.Consumer
	deviceRepo DeviceStore
	rules      *RuleEngine
	log        *slog.Logger
}

func NewAlerterService(p *kafka.Producer, c *kafka.Consumer, repo DeviceStore, rules *RuleEngine, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	if rules == nil {
		rules = NewRuleEngine()
	}

	return &Service{
		p:          p,
		c:          c,
		deviceRepo: repo,
		rules:      rules,
		log:        log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if s.c == nil {
		return fmt.Errorf("kafka consumer is not configured")
	}
	if s.p == nil {
		return fmt.Errorf("kafka producer is not configured")
	}
	if s.deviceRepo == nil {
		return fmt.Errorf("device repository is not configured")
	}

	s.log.Info("starting alerter service")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := s.c.Consume(ctx)
			if err != nil {
				s.log.Error("failed to consume processed telemetry", "err", err)
				continue
			}

			var data types.ProcessedData
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				s.log.Error("failed to unmarshal processed telemetry", "err", err)
				continue
			}

			deviceID := strconv.FormatInt(data.SenderId, 10)
			device, err := s.deviceRepo.GetDevice(ctx, deviceID)
			if err != nil {
				s.log.Error("failed to resolve device metadata", "sender_id", data.SenderId, "err", err)
				continue
			}

			alerts := s.rules.Evaluate(&data, device)
			for _, alert := range alerts {
				if err := s.p.Produce(ctx, alert); err != nil {
					s.log.Error("failed to publish alert", "device_id", alert.DeviceId, "err", err)
					continue
				}
				s.log.Info("alert published", "device_id", alert.DeviceId, "severity", alert.Severity)
			}
		}
	}
}
