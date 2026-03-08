proto:
	@protoc \
		--proto_path=api "api/telemetry.proto" \
		--go_out="internal/scrapper/gen" --go_opt=paths=source_relative \
    	--go-grpc_out="internal/scrapper/gen" --go-grpc_opt=paths=source_relative
