package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	"github.com/MotyaSS/IoTMonitoring/internal/storage/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "storage")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/storage.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error("failed to load config",
			"path", configPath,
			"err", err,
		)
		return
	}

	if cfg.Kafka == nil || cfg.Postgres == nil || cfg.Mongo == nil || cfg.S3 == nil || cfg.GRPC == nil {
		log.Error("storage config is incomplete: kafka, postgres, mongo, s3, grpc are required")
		return
	}

	consumer, err := kafka.NewConsumer(cfg.Kafka, log)
	if err != nil {
		log.Error("failed to initialize kafka consumer", "err", err)
		return
	}
	defer func() {
		_ = consumer.Close()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgresDB, err := initPostgres(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize postgres",
			"err", err,
		)
		return
	}
	defer func() {
		_ = postgresDB.Close()
	}()

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Error("failed to initialize mongo",
			"err", err,
		)
		return
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()

	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Error("failed to ping mongo",
			"err", err,
		)
		return
	}

	minioClient, err := initMinio(cfg.S3)
	if err != nil {
		log.Error("failed to initialize minio",
			"err", err,
		)
		return
	}

	metricsColl := mongoClient.Database(cfg.Mongo.DB).Collection("metrics")

	svc := service.NewStorageService(
		consumer,
		postgresDB,
		metricsColl,
		minioClient,
		cfg.S3.BucketName,
		cfg.GRPC.Addr,
		30*24*time.Hour,
		time.Hour,
		1000,
		log,
	)

	if err := svc.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("storage exited with error", "err", err)
	}

	log.Info("storage shutdown complete")

}

func initPostgres(ctx context.Context, cfg *config.PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DB,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func initMinio(cfg *config.S3Config) (*minio.Client, error) {
	secure := cfg.UseSSL
	endpoint := cfg.Endpoint

	if parsed, err := url.Parse(cfg.Endpoint); err == nil && parsed.Host != "" {
		endpoint = parsed.Host
		if parsed.Scheme == "https" {
			secure = true
		}
	}

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: secure,
	})
}
