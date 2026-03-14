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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	storagepb.UnimplementedStorageServer

	query    *queryAPI
	ingest   *ingestWorker
	archive  *archiveWorker
	grpcAddr string
	log      *slog.Logger
}

type queryAPI struct {
	postgres *sql.DB
	log      *slog.Logger
}

type ingestWorker struct {
	consumer    *kafka.Consumer
	metricsColl *mongo.Collection
	log         *slog.Logger
}

type archiveWorker struct {
	metricsColl      *mongo.Collection
	s3               *minio.Client
	bucketName       string
	retention        time.Duration
	archiveInterval  time.Duration
	archiveBatchSize int
	log              *slog.Logger
}

func NewStorageService(
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
		query: &queryAPI{
			postgres: postgres,
			log:      log,
		},
		ingest: &ingestWorker{
			consumer:    consumer,
			metricsColl: metricsColl,
			log:         log,
		},
		archive: &archiveWorker{
			metricsColl:      metricsColl,
			s3:               s3,
			bucketName:       bucketName,
			retention:        retention,
			archiveInterval:  archiveInterval,
			archiveBatchSize: archiveBatchSize,
			log:              log,
		},
		grpcAddr: grpcAddr,
		log:      log,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if s.query == nil || s.ingest == nil || s.archive == nil || s.query.postgres == nil || s.ingest.consumer == nil || s.ingest.metricsColl == nil || s.archive.s3 == nil {
		return fmt.Errorf("storage dependencies are not fully configured")
	}

	if err := s.ensureMongoIndexes(ctx); err != nil {
		return err
	}
	if err := s.ensureBucket(ctx); err != nil {
		return err
	}
	eg, childCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.runGRPC(childCtx)
	})

	eg.Go(func() error {
		return s.runArchiver(childCtx)
	})

	eg.Go(func() error {
		return s.runIngest(childCtx)
	})

	return eg.Wait()
}
