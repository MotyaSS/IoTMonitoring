package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	alerterservice "github.com/MotyaSS/IoTMonitoring/internal/alerter/service"
	alerterstorage "github.com/MotyaSS/IoTMonitoring/internal/alerter/storage"
	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	storagepb "github.com/MotyaSS/IoTMonitoring/internal/storage/gen"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "alerter")

	cfg, err := config.Load("configs/alerter.yaml")
	if err != nil || cfg.GRPC == nil || cfg.Kafka == nil || cfg.Redis == nil {
		fmt.Println("failed to load config:", err)
		return
	}

	p, err := kafka.NewProducer(cfg.Kafka, log)
	if err != nil {
		fmt.Println("failed to initialize kafka producer:", err)
		return
	}
	defer func() {
		_ = p.Close()
	}()

	c, err := kafka.NewConsumer(cfg.Kafka, log)
	if err != nil {
		fmt.Println("failed to initialize kafka consumer:", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	grpcConn, err := grpc.NewClient(cfg.GRPC.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect to storage grpc:", err)
		return
	}
	defer func() {
		_ = grpcConn.Close()
	}()

	redisDB, err := strconv.Atoi(cfg.Redis.DB)
	if err != nil {
		fmt.Println("invalid redis db value:", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       redisDB,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("failed to connect to redis:", err)
		return
	}
	defer func() {
		_ = rdb.Close()
	}()

	ttl := time.Duration(cfg.Redis.TTLSeconds) * time.Second
	deviceStore := alerterstorage.NewDeviceStore(rdb, storagepb.NewStorageClient(grpcConn), ttl, log)
	svc := alerterservice.NewAlerterService(p, c, deviceStore, alerterservice.NewRuleEngine(), log)

	if err := svc.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println("alerter exited with error:", err)
	}
	fmt.Println("alerter shutdown complete")

}
