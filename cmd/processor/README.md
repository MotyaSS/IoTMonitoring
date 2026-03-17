# processor

Сервис читает сырую телеметрию из Kafka (`raw_telemetry`), обрабатывает/нормализует данные и отправляет результат в Kafka (`processed_telemetry`).

Запуск:
```bash
CONFIG_PATH=configs/processor.yaml go run cmd/processor/main.go
```

