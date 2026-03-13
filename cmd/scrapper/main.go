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
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/server"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/service"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "scrapper")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/scrapper.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error("failed to load config", "path", configPath, "err", err)
		return
	}
	if cfg.GRPC == nil {
		log.Error("grpc missing in config file")
		return
	}
	if cfg.Kafka == nil {
		log.Error("kafka missing in config file")
		return
	}

	p, err := kafka.NewProducer(cfg.Kafka, log)
	if err != nil {
		log.Error("failed to create kafka producer", "err", err)
		return
	}
	defer func() {
		_ = p.Close()
	}()

	svc := service.NewScrapperService(p, log)

	s := server.NewServer(cfg.GRPC.Addr, svc, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("scrapper exited with error", "err", err)
	}
	log.Info("scrapper shutdown complete")
}
