package main

import (
	"context"
	"log"
	"time"

	IoTMonitoring "github.com/MotyaSS/IoTMonitoring/internal/scrapper/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const reqCount = 5

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	for range reqCount {
		<-time.After(1 * time.Second)
		cl := IoTMonitoring.NewScrapperClient(conn)
		_, err = cl.SendTelemetry(ctx, &IoTMonitoring.Telemetry{
			SenderId:    42,
			AuthToken:   "mockAuthToken",
			Latitude:    42.0,
			Longitude:   42.0,
			Timestamp:   time.Now().String(),
			Temperature: new(float32(42)),
			Pressure:    nil,
			LogMessage:  new(string("sent telemetry")),
		})
	}
}
