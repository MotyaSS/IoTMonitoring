package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	"github.com/MotyaSS/IoTMonitoring/internal/processor/service"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "processor")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/processor.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error("failed to load config", "path", configPath, "err", err)
		return
	}

	if cfg.Kafka == nil {
		log.Error("kafka missing in config file")
		return
	}

	p, err := kafka.NewProducer(cfg.Kafka, log)
	if err != nil {
		log.Error("failed to initialize kafka producer", "err", err)
		return
	}

	defer func() {
		_ = p.Close()
	}()

	c, err := kafka.NewConsumer(cfg.Kafka, log)
	if err != nil {
		log.Error("failed to initialize kafka consumer", "err", err)
		return
	}

	defer func() {
		_ = c.Close()
	}()

	s := service.NewProcessorService(p, c, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("processor exited with error", "err", err)
	}
	log.Info("processor shutdown complete")
}
