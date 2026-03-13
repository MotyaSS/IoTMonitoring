package service

import (
	"testing"
	"time"

	pb "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
)

func TestProcessData_MapsFields(t *testing.T) {
	ts := time.Now().UTC().Truncate(time.Second)
	temp := float32(21.5)

	in := &pb.Telemetry{
		SenderId:    42,
		Timestamp:   ts.Format(time.RFC3339),
		Latitude:    55.75,
		Longitude:   37.61,
		Temperature: &temp,
	}

	out, err := ProcessData(in)
	if err != nil {
		t.Fatalf("ProcessData returned error: %v", err)
	}
	if out.SenderId != 42 {
		t.Fatalf("expected sender_id 42, got %d", out.SenderId)
	}
	if !out.Timestamp.Equal(ts) {
		t.Fatalf("expected timestamp %s, got %s", ts, out.Timestamp)
	}
	if out.Latitude != 55.75 || out.Longitude != 37.61 {
		t.Fatalf("unexpected coordinates: %f %f", out.Latitude, out.Longitude)
	}
	if out.Temperature == nil || *out.Temperature != temp {
		t.Fatalf("temperature not mapped")
	}
}

func TestProcessData_InvalidTimestamp(t *testing.T) {
	in := &pb.Telemetry{Timestamp: "bad-time"}
	if _, err := ProcessData(in); err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}
