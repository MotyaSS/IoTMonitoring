package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MotyaSS/IoTMonitoring/internal/config"
	IoTMonitoring "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const reqCount = 5

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	cfg, err := config.Load("configs/scrapper.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	if cfg.GRPC == nil {
		fmt.Println("No gRPC configuration")
		return
	}

	conn, err := grpc.NewClient(cfg.GRPC.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		fmt.Println("failed to dial:", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	cl := IoTMonitoring.NewScrapperClient(conn)

	for range reqCount {
		<-time.After(1 * time.Second)
		_, err = cl.SendTelemetry(ctx, &IoTMonitoring.Telemetry{
			SenderId:    42,
			AuthToken:   "mockAuthToken",
			Latitude:    42.0,
			Longitude:   42.0,
			Timestamp:   time.Now().Format(time.RFC3339),
			Temperature: new(float32(42)),
			Pressure:    nil,
			LogMessage:  new("sent telemetry"),
		})
	}
}
