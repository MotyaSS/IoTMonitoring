package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/kafka"
	storagepb "github.com/MotyaSS/IoTMonitoring/internal/storage/gen"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	storagepb.UnimplementedStorageServer

	consumer         *kafka.Consumer
	postgres         *sql.DB
	metricsColl      *mongo.Collection
	s3               *minio.Client
	bucketName       string
	grpcAddr         string
	retention        time.Duration
	archiveInterval  time.Duration
	archiveBatchSize int
	log              *slog.Logger
}

func NewService(
	consumer *kafka.Consumer,
	postgres *sql.DB,
	metricsColl *mongo.Collection,
	s3 *minio.Client,
	bucketName, grpcAddr string,
	retention, archiveInterval time.Duration,
	archiveBatchSize int,
	log *slog.Logger,
) *Service {
	if retention <= 0 {
		retention = 30 * 24 * time.Hour
	}
	if archiveInterval <= 0 {
		archiveInterval = time.Hour
	}
	if archiveBatchSize <= 0 {
		archiveBatchSize = 1000
	}
	if log == nil {
		log = slog.Default()
	}

	return &Service{
		consumer:         consumer,
		postgres:         postgres,
		metricsColl:      metricsColl,
		s3:               s3,
		bucketName:       bucketName,
		grpcAddr:         grpcAddr,
		retention:        retention,
		archiveInterval:  archiveInterval,
		archiveBatchSize: archiveBatchSize,
		log:              log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if s.consumer == nil || s.postgres == nil || s.metricsColl == nil || s.s3 == nil {
		return fmt.Errorf("storage dependencies are not fully configured")
	}

	if err := s.ensureMongoIndexes(ctx); err != nil {
		return err
	}
	if err := s.ensureBucket(ctx); err != nil {
		return err
	}

	errCh := make(chan error, 3)

	go func() {
		if err := s.runGRPC(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := s.runArchiver(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := s.runIngest(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
