# scrapper

Сервис принимает телеметрию по gRPC и отправляет сообщения в Kafka (`raw_telemetry`).

Запуск:
```bash
CONFIG_PATH=configs/scrapper.yaml go run cmd/scrapper/main.go
```

