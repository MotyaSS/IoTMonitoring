package main

import (
	"context"
	"errors"
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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/alerter.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error("failed to load config", "path", configPath, "err", err)
		return
	}
	if cfg.GRPC == nil || cfg.Kafka == nil || cfg.Redis == nil {
		log.Error("alerter config is incomplete: grpc, kafka and redis are required")
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

	grpcConn, err := grpc.NewClient(cfg.GRPC.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to storage grpc", "addr", cfg.GRPC.Addr, "err", err)
		return
	}

	defer func() {
		_ = grpcConn.Close()
	}()

	redisDB, err := strconv.Atoi(cfg.Redis.DB)
	if err != nil {
		log.Error("invalid redis db value", "db", cfg.Redis.DB, "err", err)
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
		log.Error("failed to connect to redis", "addr", cfg.Redis.Addr, "err", err)
		return
	}

	defer func() {
		_ = rdb.Close()
	}()

	ttl := time.Duration(cfg.Redis.TTLSeconds) * time.Second
	deviceStore := alerterstorage.NewDeviceStore(rdb, storagepb.NewStorageClient(grpcConn), ttl, log)
	svc := alerterservice.NewAlerterService(p, c, deviceStore, alerterservice.NewRuleEngine(), log)

	if err := svc.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("alerter exited with error", "err", err)
	}
	log.Info("alerter shutdown complete")

}
