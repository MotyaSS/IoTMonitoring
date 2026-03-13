# Storage service

Service responsibilities:

- Read processed telemetry from Kafka topic `processed_telemetry`.
- Store incoming metrics in Mongo collection `metrics`.
- Keep hot metrics in Mongo for 30 days.
- Archive metrics older than 30 days to Minio (`jsonl.gz`) and then delete them from Mongo.
- Serve gRPC method `GetDevice` for `alerter` by reading thresholds from Postgres table `devices`.

Archive object key format:

`metrics/year=YYYY/month=MM/day=DD/part-<first_id>-<last_id>.jsonl.gz`
