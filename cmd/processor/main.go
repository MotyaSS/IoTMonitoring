package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	"github.com/MotyaSS/IoTMonitoring/internal/processor/kafka/consumer"
	"github.com/MotyaSS/IoTMonitoring/internal/processor/kafka/producer"
	"github.com/MotyaSS/IoTMonitoring/internal/processor/service"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "processor")

	cfg, err := config.Load("configs/processor.yaml")
	if err != nil {
		fmt.Println("failed to load config: ", err)
		return
	}

	if cfg.Kafka == nil {
		fmt.Println("kafka missing in config file")
		return
	}

	p, err := producer.NewProcessorProducer(cfg.Kafka, log)
	if err != nil {
		fmt.Println("failed to initialize kafka producer: ", err)
		return
	}

	defer func() {
		_ = p.Close()
	}()

	c, err := consumer.NewProcessorConsumer(cfg.Kafka, log)
	if err != nil {
		fmt.Println("failed to initialize kafka consumer: ", err)
		return
	}

	defer func() {
		_ = c.Close()
	}()

	s := service.NewProcessorService(p, c, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println("processor exited with error:", err)
	}
	fmt.Println("processor shutdown complete")
}
