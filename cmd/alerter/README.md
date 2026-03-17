# alerter

Сервис читает `processed_telemetry` из Kafka, проверяет метрики по порогам устройства и публикует алерты в Kafka (`alert`).

Использует:
- gRPC к `storage` для получения данных устройства (если нет в кэше)
- Redis для кэширования порогов устройств

Запуск:
```bash
CONFIG_PATH=configs/alerter.yaml go run cmd/alerter/main.go
```

