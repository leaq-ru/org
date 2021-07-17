FROM golang:1.15-alpine AS build
WORKDIR /app
COPY / /app
RUN go build -o servicebin

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/app/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /app/grpc_health_probe

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/servicebin /app
COPY --from=build /app/grpc_health_probe /app
