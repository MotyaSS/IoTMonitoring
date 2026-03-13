package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Service) runArchiver(ctx context.Context) error {
	ticker := time.NewTicker(s.archiveInterval)
	defer ticker.Stop()

	s.log.Info("storage archiver started", "interval", s.archiveInterval.String(), "retention", s.retention.String())

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.archiveOnce(ctx); err != nil {
				s.log.Error("archive cycle failed", "err", err)
			}
		}
	}
}

func (s *Service) archiveOnce(ctx context.Context) error {
	cutoff := time.Now().UTC().Add(-s.retention)

	filter := bson.M{"timestamp": bson.M{"$lt": cutoff}}
	findOpts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: 1}, {Key: "_id", Value: 1}}).
		SetLimit(int64(s.archiveBatchSize))

	cur, err := s.metricsColl.Find(ctx, filter, findOpts)
	if err != nil {
		return fmt.Errorf("find metrics to archive: %w", err)
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	docs := make([]metricDoc, 0, s.archiveBatchSize)
	for cur.Next(ctx) {
		var d metricDoc
		if err := cur.Decode(&d); err != nil {
			return fmt.Errorf("decode archive document: %w", err)
		}
		docs = append(docs, d)
	}
	if err := cur.Err(); err != nil {
		return fmt.Errorf("iterate archive cursor: %w", err)
	}
	if len(docs) == 0 {
		return nil
	}

	objectKey, payload, err := buildArchiveObject(docs)
	if err != nil {
		return fmt.Errorf("build archive object: %w", err)
	}

	_, err = s.s3.PutObject(ctx, s.bucketName, objectKey, bytes.NewReader(payload), int64(len(payload)), minio.PutObjectOptions{
		ContentType: "application/gzip",
	})
	if err != nil {
		return fmt.Errorf("put object to minio: %w", err)
	}

	ids := make([]primitive.ObjectID, 0, len(docs))
	for _, d := range docs {
		ids = append(ids, d.ID)
	}

	if _, err := s.metricsColl.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}}); err != nil {
		return fmt.Errorf("delete archived metrics from mongo: %w", err)
	}

	s.log.Info("archive batch completed", "count", len(docs), "object", objectKey)
	return nil
}

func buildArchiveObject(docs []metricDoc) (string, []byte, error) {
	if len(docs) == 0 {
		return "", nil, fmt.Errorf("no documents to archive")
	}

	first := docs[0]
	last := docs[len(docs)-1]
	t := first.Timestamp.UTC()
	objectKey := fmt.Sprintf(
		"metrics/year=%04d/month=%02d/day=%02d/part-%s-%s.jsonl.gz",
		t.Year(), int(t.Month()), t.Day(), first.ID.Hex(), last.ID.Hex(),
	)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	enc := json.NewEncoder(gz)

	for _, d := range docs {
		if err := enc.Encode(d); err != nil {
			_ = gz.Close()
			return "", nil, err
		}
	}

	if err := gz.Close(); err != nil {
		return "", nil, err
	}

	payload, err := io.ReadAll(&buf)
	if err != nil {
		return "", nil, err
	}

	return objectKey, payload, nil
}
