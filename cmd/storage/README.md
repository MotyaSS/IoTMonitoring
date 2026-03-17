# storage

Сервис читает `processed_telemetry` из Kafka и:
- сохраняет свежие метрики в MongoDB
- хранит данные устройств в PostgreSQL (таблица `devices`)
- архивирует метрики старше 30 дней в MinIO
- поднимает gRPC API для `alerter` (получение данных устройства)

Запуск:
```bash
CONFIG_PATH=configs/storage.yaml go run cmd/storage/main.go
```

