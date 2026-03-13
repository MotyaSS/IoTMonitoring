package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	storagepb "github.com/MotyaSS/IoTMonitoring/internal/storage/gen"
	"github.com/MotyaSS/IoTMonitoring/internal/types"
	"github.com/redis/go-redis/v9"
)

type DeviceStore struct {
	rdb    *redis.Client
	grpc   storagepb.StorageClient
	log    *slog.Logger
	ttl    time.Duration
	prefix string
}

func NewDeviceStore(rdb *redis.Client, grpcClient storagepb.StorageClient, ttl time.Duration, log *slog.Logger) *DeviceStore {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	if log == nil {
		log = slog.Default()
	}

	return &DeviceStore{
		rdb:    rdb,
		grpc:   grpcClient,
		log:    log,
		ttl:    ttl,
		prefix: "device:",
	}
}

func (s *DeviceStore) GetDevice(ctx context.Context, deviceID string) (*types.Device, error) {
	if deviceID == "" {
		return nil, errors.New("device id is empty")
	}

	if s.rdb != nil {
		if d, err := s.getFromCache(ctx, deviceID); err == nil {
			return d, nil
		} else if !errors.Is(err, redis.Nil) {
			s.log.Warn("failed to load device from redis", "device_id", deviceID, "err", err)
		}
	}

	if s.grpc == nil {
		return nil, errors.New("storage grpc client is not configured")
	}

	resp, err := s.grpc.GetDevice(ctx, &storagepb.GetDeviceRequest{DeviceId: deviceID})
	if err != nil {
		return nil, fmt.Errorf("storage grpc get device: %w", err)
	}
	if resp.GetDevice() == nil {
		return nil, fmt.Errorf("device %s not found", deviceID)
	}

	device := mapPBDevice(resp.GetDevice())
	if s.rdb != nil {
		s.saveToCache(ctx, device)
	}

	return device, nil
}

func (s *DeviceStore) getFromCache(ctx context.Context, deviceID string) (*types.Device, error) {
	key := s.prefix + deviceID
	raw, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var d types.Device
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *DeviceStore) saveToCache(ctx context.Context, d *types.Device) {
	if d == nil {
		return
	}

	payload, err := json.Marshal(d)
	if err != nil {
		s.log.Warn("failed to marshal device for redis", "device_id", d.DeviceId, "err", err)
		return
	}

	if err := s.rdb.Set(ctx, s.prefix+d.DeviceId, payload, s.ttl).Err(); err != nil {
		s.log.Warn("failed to cache device in redis", "device_id", d.DeviceId, "err", err)
	}
}

func mapPBDevice(d *storagepb.Device) *types.Device {
	return &types.Device{
		DeviceId:             d.GetDeviceId(),
		WindSpeedHighBound:   d.GetWindSpeedHighBound(),
		WindSpeedLowBound:    d.GetWindSpeedLowBound(),
		TemperatureHighBound: d.GetTemperatureHighBound(),
		TemperatureLowBound:  d.GetTemperatureLowBound(),
		PressureHighBound:    d.GetPressureHighBound(),
		PressureLowBound:     d.GetPressureLowBound(),
	}
}
