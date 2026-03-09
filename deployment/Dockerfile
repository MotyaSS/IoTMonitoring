# Multi-stage Dockerfile to build each command binary by passing --build-arg SERVICE=<name>
ARG SERVICE=scrapper

FROM golang:1.21-alpine AS builder
ARG SERVICE
WORKDIR /src

# Install build deps
RUN apk add --no-cache git ca-certificates

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Build the specific service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o /app/bin/${SERVICE} ./cmd/${SERVICE}

# Final image
FROM alpine:3.18
ARG SERVICE
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/${SERVICE} /app/${SERVICE}

# default config path (can be overridden with a mounted file or env)
ENV CONFIG_PATH=/app/config.yaml

# expose a common gRPC port (services can choose different ports in their config)
EXPOSE 8080

ENTRYPOINT ["/app/${SERVICE}"]

