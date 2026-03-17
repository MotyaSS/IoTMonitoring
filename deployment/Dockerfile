# you can change SERVICE value to build different service
ARG SERVICE=scrapper
ARG GO_VERSION=1.26.1

FROM golang:${GO_VERSION}-alpine AS builder
ARG SERVICE
WORKDIR /src

# Install build deps
RUN apk add --no-cache git ca-certificates

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Build the selected service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o /app/bin/${SERVICE} ./cmd/${SERVICE}

# Final image
FROM alpine:3.18
ARG SERVICE
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/${SERVICE} /app/service

# default config path (can be overridden with a mounted file or env)
ENV CONFIG_PATH=/app/config.yaml

# Expose a common gRPC port (services may use different ports via config)
EXPOSE 8080

ENTRYPOINT ["/app/service"]

