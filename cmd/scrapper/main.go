package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/kafka/producer"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/server"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/service"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "scrapper")
	cfg, err := config.Load("configs/scrapper.yaml")
	if err != nil {
		fmt.Println("failed to load config: ", err)
		return
	}
	if cfg.GRPC == nil {
		fmt.Println("grpc missing in config file")
		return
	}
	if cfg.Kafka == nil {
		fmt.Println("kafka missing in config file")
		return
	}

	kp, err := producer.NewKafkaProducer(cfg.Kafka)
	if err != nil {
		fmt.Println("failed to create kafka producer: ", err)
	}
	defer kp.Close()

	svc := service.New(kp)

	s := server.NewServer(cfg.GRPC.Addr, svc, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := s.Run(ctx); err != nil {
		fmt.Println("server exited with error:", err)
	}

	fmt.Println("server shutdown complete")
}
