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
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: 1}},
		},
	}

	_, err := s.metricsColl.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create mongo indexes: %w", err)
	}
	return nil
}

func (s *Service) ensureBucket(ctx context.Context) error {
	exists, err := s.s3.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("check minio bucket: %w", err)
	}
	if exists {
		return nil
	}
	if err := s.s3.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create minio bucket: %w", err)
	}
	return nil
}
