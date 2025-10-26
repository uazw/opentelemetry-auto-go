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
   curl http://localhost:8080/hello
   ```
3. Observe telemetry:
   - Collector logs (metrics and traces): `docker compose logs -f otel-collector`
   - App Prometheus metrics endpoint (scraped by the collector): `http://localhost:8080/metrics`

## Project Structure
- `cmd/server/main.go` – GoFrame web server with `/` and `/hello`
- `Dockerfile` – Multi-stage build for the app image
- `docker-compose.yml` – Brings up app, OTEL Collector, and Go auto-instrumentation agent
- `otel-collector-config.yaml` – Collector configuration (OTLP receiver, debug + Prometheus exporters, traces pipeline)

## Configuration
- Auto-instrumentation agent environment (see `docker-compose.yml` service `otel-go-agent`):
  - `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317`
  - `OTEL_EXPORTER_OTLP_PROTOCOL=grpc`
  - `OTEL_SERVICE_NAME=otel-go-webapp`
  - `OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev`
  - `OTEL_GO_AUTO_INSTRUMENTATION_HTTP_ENABLED=true`
  - `OTEL_GO_AUTO_INSTRUMENTATION_RUNTIME_ENABLED=true`
  - `OTEL_GO_AUTO_TARGET_EXE=/usr/local/bin/server` (target binary inside the app container)

## Development
The app can run standalone without telemetry: `go run ./cmd/server`. Auto-instrumentation requires Linux with eBPF support and container privileges; use Docker Compose to enable it.

## Notes
- The auto-instrumentation agent container joins the app container's PID namespace to attach and collect telemetry.
- This setup typically requires Linux; Docker Desktop on macOS/Windows may not support the required kernel features.
- Pin images for reproducibility if desired (collector and agent).

## Excluding /metrics Telemetry (Collector)
This repo includes an OTTL filter in the OpenTelemetry Collector to drop telemetry related to the app's own `/metrics` endpoint.

- Location: `otel-collector-config.yaml` under `processors.filter/ottl` and wired into both `metrics` and `traces` pipelines.
- Behavior:
  - Drops metrics datapoints where `attributes["http_route"] == "/metrics"`.
  - Drops spans where `attributes["url.path"] == "/metrics"`.
- Pipelines apply `processors: [filter/ottl, batch]` to enforce the filter.

To modify or disable the filter, edit `otel-collector-config.yaml` and restart the collector:

```
docker compose restart otel-collector
```


## Related Doc
[auto-go](https://github.com/open-telemetry/opentelemetry-go-instrumentation/blob/main/docs/getting-started.md)
[how-to-mix-manual-and-auto](https://goframe.org/en/docs/obs/metrics-builtin)
