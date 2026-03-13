package service

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (s *Service) ensureMongoIndexes(ctx context.Context) error {
	if s.ingest == nil || s.ingest.metricsColl == nil {
		return fmt.Errorf("ingest metrics collection is not configured")
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: 1}},
		},
	}

	_, err := s.ingest.metricsColl.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create mongo indexes: %w", err)
	}
	return nil
}

func (s *Service) ensureBucket(ctx context.Context) error {
	if s.archive == nil || s.archive.s3 == nil {
		return fmt.Errorf("archive s3 client is not configured")
	}

	exists, err := s.archive.s3.BucketExists(ctx, s.archive.bucketName)
	if err != nil {
		return fmt.Errorf("check minio bucket: %w", err)
	}
	if exists {
		return nil
	}
	if err := s.archive.s3.MakeBucket(ctx, s.archive.bucketName, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create minio bucket: %w", err)
	}
	return nil
}
