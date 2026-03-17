# IoTMonitoring

Небольшая микросервисная система для обработки IoT-телеметрии:
- `scrapper` принимает сырые данные по gRPC и пишет в Kafka (`raw_telemetry`)
- `processor` читает `raw_telemetry`, нормализует данные и пишет в `processed_telemetry`
- `storage` читает `processed_telemetry`, хранит свежие метрики в MongoDB, старые архивирует в MinIO, и отдает `devices` по gRPC
- `alerter` читает `processed_telemetry`, подтягивает пороги устройства из `storage` по gRPC (с Redis-кэшем), публикует алерты в Kafka (`alert`)

## Что нужно для запуска

- Go `1.26.1` (см. `go.mod`)
- Docker + Docker Compose (для полного запуска стека)
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` (если нужно пересобирать protobuf)

---
## Быстрый запуск в контейнерах

#### Поднять:
```bash
docker compose -f deployment/docker-compose.yml up -d --build
```
или 
```bash
make docker-up
```

#### Посмотреть логи:

```bash
docker compose -f deployment/docker-compose.yml ps
docker compose -f deployment/docker-compose.yml logs -f --tail=150
```
или
```bash
make docker-logs
```

#### Остановить:

```bash
docker compose -f deployment/docker-compose.yml down
```
или
```bash
make docker-down
```

> Если поднимали проект раньше и меняли init-скрипты Postgres, для чистой инициализации удалите volume:
>
> ```bash
> docker compose -f deployment/docker-compose.yml down -v
> ```

---
## Локальный запуск

Конфиги из `configs/*.yaml` по умолчанию ориентированы на docker-сеть (`kafka`, `postgres`, `mongo`, `redis`, `minio` как hostnames).
Для запуска на хосте обновите адреса в конфигах на `localhost`.

#### Пример команд:

```bash
CONFIG_PATH=configs/scrapper.yaml go run cmd/scrapper/main.go
CONFIG_PATH=configs/processor.yaml go run cmd/processor/main.go
CONFIG_PATH=configs/storage.yaml go run cmd/storage/main.go
CONFIG_PATH=configs/alerter.yaml go run cmd/alerter/main.go
```

---
## Миграции / инициализация БД

- SQL-инициализация таблицы `devices` лежит в `deployment/postgres/init/001_devices.sql`
- Скрипт автоматически применяется контейнером Postgres при первом старте пустого volume

## Полезные команды

```bash
go test ./...
make proto
make docker-up
make docker-down
```

