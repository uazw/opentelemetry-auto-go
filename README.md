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
   - Auto-instrumentation agent logs: `docker compose logs -f otel-go-agent`
   - App Prometheus metrics endpoint (scraped by the collector): `http://localhost:8080/metrics`

## Run Locally (no Collector/Agent)
- Start the app directly: `go run ./cmd/server`
- Visit: `http://localhost:8080/`, `http://localhost:8080/hello`, `http://localhost:8080/metrics`

## Project Structure
- `cmd/server/main.go` – GoFrame web server with `/`, `/hello`, and `/metrics`
- `cmd/server/JsonOutputsForLogger.go` – Custom JSON logger handler for structured logs
- `Dockerfile` – Multi-stage build for the app image
- `docker-compose.yml` – Brings up app, OTEL Collector, and Go auto-instrumentation agent
- `otel-collector-config.yaml` – Collector configuration (OTLP receiver, debug exporter, Prometheus scrape, OTTL filters)

## Manual Metrics (in App)
This repo includes manual application metrics in addition to auto-instrumentation:

- Prometheus exporter and OTEL metrics provider are initialized at startup so the app can expose metrics directly at `http://localhost:8080/metrics`.
- A demo counter metric is created and incremented once at startup.
- The `/metrics` endpoint is served by GoFrame's OTEL Prometheus handler and is scraped by the OpenTelemetry Collector.

Key locations in code:

- JSON logger handler set: `cmd/server/main.go:23`
- Prometheus exporter creation: `cmd/server/main.go:25`
- OTEL metrics provider setup: `cmd/server/main.go:41`–`cmd/server/main.go:46`
- Counter creation: `cmd/server/main.go:32`–`cmd/server/main.go:38`
- HTTP handlers (`/`, `/hello`): `cmd/server/main.go:52`–`cmd/server/main.go:60`
- Metrics endpoint binding: `cmd/server/main.go:61`
- Counter increment: `cmd/server/main.go:62`

Extending manual metrics:

- Increment the counter per request (e.g., inside each handler) or add histograms for latency and request size.
- You can add attributes/labels by using metric instruments from the OTEL metrics API via GoFrame's `gmetric` provider.

Note: Manual tracing is not implemented in this app. Traces observed in the Collector come from the Go auto-instrumentation agent. If you want manual spans, wire an OTLP trace exporter and use `otel.Tracer` in handlers.

## Configuration
- Auto-instrumentation agent environment (see `docker-compose.yml` service `otel-go-agent`):
  - `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317`
  - `OTEL_EXPORTER_OTLP_PROTOCOL=grpc`
  - `OTEL_SERVICE_NAME=otel-go-webapp`
  - `OTEL_RESOURCE_ATTRIBUTES=deployment.environment=dev`
  - `OTEL_GO_AUTO_INSTRUMENTATION_HTTP_ENABLED=true`
  - `OTEL_GO_AUTO_INSTRUMENTATION_RUNTIME_ENABLED=true`
  - `OTEL_GO_AUTO_TARGET_EXE=/usr/local/bin/server` (target binary inside the app container)
  - `OTEL_LOG_LEVEL=debug`

## Development
The app can run standalone without telemetry: `go run ./cmd/server`. Auto-instrumentation requires Linux with eBPF support and container privileges; use Docker Compose to enable it.

## Notes
- The auto-instrumentation agent container joins the app container's PID namespace to attach and collect telemetry.
- This setup typically requires Linux; Docker Desktop on macOS/Windows may not support the required kernel features.
- Host networking is used for simplicity; if not suitable in your environment, remove `network_mode: "host"` and rely on published ports, adjusting endpoints accordingly.
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

## Troubleshooting
- Agent cannot attach / no traces:
  - Ensure you run on Linux with eBPF support and sufficient privileges. Check agent logs: `docker compose logs -f otel-go-agent`.
  - Confirm the app binary path matches `OTEL_GO_AUTO_TARGET_EXE=/usr/local/bin/server` inside the app container.
- Collector not receiving data:
  - Check collector logs: `docker compose logs -f otel-collector`.
  - Verify endpoints/protocols (`4317` gRPC, `4318` HTTP) and that `OTEL_EXPORTER_OTLP_ENDPOINT` matches.
- Metrics missing in collector output:
  - Ensure the app `/metrics` endpoint is reachable from the collector (host networking assumed); open `http://localhost:8080/metrics`.
  - Remember the OTTL filter drops `/metrics`-route telemetry from the pipelines by design.
