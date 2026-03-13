package service

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMakeEventIDDeterministic(t *testing.T) {
	payload := []byte(`{"sender_id":1,"temperature":10.5}`)
	first := makeEventID(payload)
	second := makeEventID(payload)

	if first == "" {
		t.Fatal("event id must not be empty")
	}
	if first != second {
		t.Fatalf("expected deterministic hash, got %s and %s", first, second)
	}
}

func TestBuildArchiveObject(t *testing.T) {
	now := time.Now().UTC()
	docs := []metricDoc{
		{ID: bson.NewObjectID(), EventID: "e1", DeviceID: "1", Timestamp: now, Latitude: 1.0, Longitude: 2.0, ReceivedAt: now},
		{ID: bson.NewObjectID(), EventID: "e2", DeviceID: "1", Timestamp: now.Add(time.Minute), Latitude: 1.1, Longitude: 2.1, ReceivedAt: now},
	}

	key, payload, err := buildArchiveObject(docs)
	if err != nil {
		t.Fatalf("buildArchiveObject failed: %v", err)
	}
	if !strings.HasPrefix(key, "metrics/year=") {
		t.Fatalf("unexpected object key: %s", key)
	}
	if len(payload) == 0 {
		t.Fatal("archive payload must not be empty")
	}

	zr, err := gzip.NewReader(bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("gzip.NewReader failed: %v", err)
	}
	defer func() {
		_ = zr.Close()
	}()

	s := bufio.NewScanner(zr)
	lines := 0
	for s.Scan() {
		lines++
		var decoded map[string]any
		if err := json.Unmarshal(s.Bytes(), &decoded); err != nil {
			t.Fatalf("line is not valid json: %v", err)
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if lines != len(docs) {
		t.Fatalf("expected %d lines, got %d", len(docs), lines)
	}
}
