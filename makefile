proto:
	@protoc \
		--proto_path=api "api/telemetry.proto" \
		--go_out="internal/scrapper/gen" --go_opt=paths=source_relative \
	    	--go-grpc_out="internal/scrapper/gen" --go-grpc_opt=paths=source_relative
	@mkdir -p internal/storage/gen
	@protoc \
		--proto_path=api "api/storage.proto" \
		--go_out="internal/storage/gen" --go_opt=paths=source_relative \
		--go-grpc_out="internal/storage/gen" --go-grpc_opt=paths=source_relative
scrapper:
	@go run cmd/scrapper/main.go
processor:
	@go run cmd/processor/main.go
alerter:
	@go run cmd/alerter/main.go
storage:
	@go run cmd/storage/main.go
client: # mock
	@go run cmd/client_mock/main.go

docker-up:
	@docker compose -f deployment/docker-compose.yml up -d --build

docker-down:
	@docker compose -f deployment/docker-compose.yml down

docker-logs:
	@docker compose -f deployment/docker-compose.yml logs -f --tail=150
