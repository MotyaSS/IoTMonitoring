package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	storagepb "github.com/MotyaSS/IoTMonitoring/internal/storage/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetDevice(ctx context.Context, req *storagepb.GetDeviceRequest) (*storagepb.GetDeviceResponse, error) {
	if s.query == nil {
		return nil, status.Error(codes.Internal, "query layer is not configured")
	}
	return s.query.GetDevice(ctx, req)
}

func (q *queryAPI) GetDevice(ctx context.Context, req *storagepb.GetDeviceRequest) (*storagepb.GetDeviceResponse, error) {
	deviceID := req.GetDeviceId()
	if deviceID == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id is required")
	}

	query := `
		SELECT device_id,
		       wind_speed_high_bound,
		       wind_speed_low_bound,
		       temperature_high_bound,
		       temperature_low_bound,
		       pressure_high_bound,
		       pressure_low_bound
		FROM devices
		WHERE device_id = $1
	`

	var d storagepb.Device
	err := q.postgres.QueryRowContext(ctx, query, deviceID).Scan(
		&d.DeviceId,
		&d.WindSpeedHighBound,
		&d.WindSpeedLowBound,
		&d.TemperatureHighBound,
		&d.TemperatureLowBound,
		&d.PressureHighBound,
		&d.PressureLowBound,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "device not found")
		}
		q.log.Error("postgres query failed",
			"device_id", deviceID,
			"err", err,
		)
		return nil, status.Error(codes.Internal, "failed to query device")
	}

	return &storagepb.GetDeviceResponse{Device: &d}, nil
}

func (s *Service) runGRPC(ctx context.Context) error {
	if s.grpcAddr == "" {
		return fmt.Errorf("grpc address is empty")
	}

	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}

	g := grpc.NewServer()
	storagepb.RegisterStorageServer(g, s)

	errCh := make(chan error, 1)
	go func() {
		errCh <- g.Serve(lis)
	}()

	s.log.Info("storage grpc server started",
		"addr", s.grpcAddr,
	)

	select {
	case <-ctx.Done():
		g.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
