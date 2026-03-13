package types

type Device struct {
	DeviceId             string  `json:"device_id" db:"device_id"`
	WindSpeedHighBound   float64 `json:"wind_speed_high_bound" db:"wind_speed_high_bound"`
	WindSpeedLowBound    float64 `json:"wind_speed_low_bound" db:"wind_speed_low_bound"`
	TemperatureHighBound float64 `json:"temperature_high_bound" db:"temperature_high_bound"`
	TemperatureLowBound  float64 `json:"temperature_low_bound" db:"temperature_low_bound"`
	PressureHighBound    float64 `json:"pressure_high_bound" db:"pressure_high_bound"`
	PressureLowBound     float64 `json:"pressure_low_bound" db:"pressure_low_bound"`
}
