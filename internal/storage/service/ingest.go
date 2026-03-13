package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *Service) runIngest(ctx context.Context) error {
	s.log.Info("storage ingest loop started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := s.consumer.Consume(ctx)
			if err != nil {
				s.log.Error("failed to consume processed telemetry", "err", err)
				continue
			}

			var data types.ProcessedData
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				s.log.Error("failed to unmarshal processed telemetry", "err", err)
				continue
			}

			doc := metricDoc{
				EventID:     makeEventID(msg.Value),
				SenderID:    data.SenderId,
				DeviceID:    strconv.FormatInt(data.SenderId, 10),
				Timestamp:   data.Timestamp.UTC(),
				Latitude:    data.Latitude,
				Longitude:   data.Longitude,
				LogMessage:  data.LogMessage,
				Temperature: data.Temperature,
				Pressure:    data.Pressure,
				WindSpeed:   data.WindSpeed,
				ReceivedAt:  time.Now().UTC(),
			}

			if err := s.insertMetric(ctx, doc); err != nil {
				s.log.Error("failed to store metric in mongo", "err", err)
				continue
			}
		}
	}
}

func (s *Service) insertMetric(ctx context.Context, doc metricDoc) error {
	_, err := s.metricsColl.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	return nil
}
