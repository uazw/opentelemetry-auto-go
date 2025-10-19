# Go Web App with OTEL Go Auto-Instrumentation

This repo contains a simple Go web server and a Docker Compose setup that uses OpenTelemetry Go Auto-Instrumentation to collect metrics and traces and send them to an OpenTelemetry Collector.

## Prerequisites
- Docker and Docker Compose installed

## Run with Docker Compose
1. Build and start services:
   ```
   docker compose up --build
   ```
2. Generate traffic:
   ```
   curl http://localhost:8080/
   ```
3. Observe telemetry:
   - Collector logs (metrics and traces): `docker compose logs -f otel-collector`
   - Prometheus metrics endpoint from collector: `http://localhost:8889/metrics`

## Project Structure
- `cmd/server/main.go` – Minimal Go web server (no manual OTEL code)
- `Dockerfile` – Multi-stage build for the app image
- `docker-compose.yml` – Brings up app, OTEL Collector, and Go auto-instrumentation agent
- `otel-collector-config.yaml` – Collector configuration (OTLP receiver, debug + Prometheus exporters, traces pipeline)

## Configuration
- Auto-instrumentation agent environment:
  - `OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317`
  - `OTEL_SERVICE_NAME=otel-go-webapp`
  - `OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev`
  - `GOAUTO_INSTRUMENTATION_HTTP_ENABLED=true`
  - `GOAUTO_INSTRUMENTATION_RUNTIME_ENABLED=true`
  - `GOAUTO_TARGET_EXE=/usr/local/bin/server` (target binary inside the app container)

## Development
The app can run standalone without telemetry: `go run ./cmd/server`. Auto-instrumentation requires Linux with eBPF support and container privileges; use Docker Compose to enable it.

## Notes
- The auto-instrumentation agent container joins the app container's PID namespace to attach and collect telemetry.
- This setup typically requires Linux; Docker Desktop on macOS/Windows may not support the required kernel features.
- Pin images for reproducibility if desired (collector and agent).
