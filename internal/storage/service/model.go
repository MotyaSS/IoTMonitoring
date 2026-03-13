package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type metricDoc struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"-"`
	EventID     string        `bson:"event_id" json:"event_id"`
	SenderID    int64         `bson:"sender_id" json:"sender_id"`
	DeviceID    string        `bson:"device_id" json:"device_id"`
	Timestamp   time.Time     `bson:"timestamp" json:"timestamp"`
	Latitude    float64       `bson:"latitude" json:"latitude"`
	Longitude   float64       `bson:"longitude" json:"longitude"`
	LogMessage  *string       `bson:"log_message,omitempty" json:"log_message,omitempty"`
	Temperature *float32      `bson:"temperature,omitempty" json:"temperature,omitempty"`
	Pressure    *float32      `bson:"pressure,omitempty" json:"pressure,omitempty"`
	WindSpeed   *float32      `bson:"wind_speed,omitempty" json:"wind_speed,omitempty"`
	ReceivedAt  time.Time     `bson:"received_at" json:"received_at"`
}

func makeEventID(payload []byte) string {
	h := sha256.Sum256(payload)
	return hex.EncodeToString(h[:])
}
