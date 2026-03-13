package service

import (
	"fmt"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/types"
)

type RuleEngine struct{}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (re *RuleEngine) Evaluate(data *types.ProcessedData, device *types.Device) []types.Alert {
	if data == nil || device == nil {
		return nil
	}

	now := time.Now().UTC()
	alerts := make([]types.Alert, 0, 3)

	if data.Temperature != nil {
		value := float64(*data.Temperature)
		if value > device.TemperatureHighBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "high", fmt.Sprintf("temperature %.2f is above high bound %.2f", value, device.TemperatureHighBound)))
		}
		if value < device.TemperatureLowBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "high", fmt.Sprintf("temperature %.2f is below low bound %.2f", value, device.TemperatureLowBound)))
		}
	}

	if data.Pressure != nil {
		value := float64(*data.Pressure)
		if value > device.PressureHighBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "medium", fmt.Sprintf("pressure %.2f is above high bound %.2f", value, device.PressureHighBound)))
		}
		if value < device.PressureLowBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "medium", fmt.Sprintf("pressure %.2f is below low bound %.2f", value, device.PressureLowBound)))
		}
	}

	if data.WindSpeed != nil {
		value := float64(*data.WindSpeed)
		if value > device.WindSpeedHighBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "medium", fmt.Sprintf("wind speed %.2f is above high bound %.2f", value, device.WindSpeedHighBound)))
		}
		if value < device.WindSpeedLowBound {
			alerts = append(alerts, newAlert(device.DeviceId, now, "medium", fmt.Sprintf("wind speed %.2f is below low bound %.2f", value, device.WindSpeedLowBound)))
		}
	}

	return alerts
}

func newAlert(deviceID string, ts time.Time, severity, log string) types.Alert {
	return types.Alert{
		AlertId:   fmt.Sprintf("%s-%d", deviceID, ts.UnixNano()),
		DeviceId:  deviceID,
		Severity:  severity,
		TimeStamp: ts,
		Log:       log,
	}
}
