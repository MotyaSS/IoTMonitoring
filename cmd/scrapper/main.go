package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MotyaSS/IoTMonitoring/internal/logger"
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/server"
)

func main() {
	log := logger.NewLoggerWithPrefix(slog.NewTextHandler(os.Stdout, nil), "scrapper")

	s := server.NewServer(":8080", log)
	fmt.Println(s.Run())
}
