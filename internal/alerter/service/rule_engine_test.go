package service

import (
	"testing"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/types"
)

func TestRuleEngineEvaluate_GeneratesAlerts(t *testing.T) {
	re := NewRuleEngine()
	temp := float32(90)
	pressure := float32(400)

	data := &types.ProcessedData{
		SenderId:    10,
		Timestamp:   time.Now(),
		Temperature: &temp,
		Pressure:    &pressure,
	}
	device := &types.Device{
		DeviceId:             "10",
		TemperatureHighBound: 50,
		TemperatureLowBound:  -20,
		PressureHighBound:    300,
		PressureLowBound:     100,
	}

	alerts := re.Evaluate(data, device)
	if len(alerts) < 2 {
		t.Fatalf("expected at least 2 alerts, got %d", len(alerts))
	}
	for _, alert := range alerts {
		if alert.DeviceId != "10" {
			t.Fatalf("unexpected device id: %s", alert.DeviceId)
		}
		if alert.AlertId == "" {
			t.Fatal("expected non-empty alert id")
		}
	}
}

func TestRuleEngineEvaluate_NilInput(t *testing.T) {
	re := NewRuleEngine()
	if alerts := re.Evaluate(nil, nil); len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}
