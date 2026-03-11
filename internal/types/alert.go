package types

import "time"

type Alert struct {
	AlertId   string    `json:"alert_id"`
	DeviceId  string    `json:"device_id"`
	Severity  string    `json:"severity"`
	TimeStamp time.Time `json:"timestamp"`
	Log       string    `json:"log"`
}
