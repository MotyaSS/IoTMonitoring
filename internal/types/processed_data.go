package types

type ProcessedData struct {
	SenderId    int64    `json:"sender_id,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
	Latitude    float64  `json:"latitude,omitempty"`
	Longitude   float64  `json:"longitude,omitempty"`
	LogMessage  *string  `json:"log_message,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	Pressure    *float32 `json:"pressure,omitempty"`
	WindSpeed   *float32 `json:"wind_speed,omitempty"`
}
