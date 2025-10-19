package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
)

func main() {
	var ctx = gctx.New()

	// Prometheus exporter to export metrics as Prometheus format.
	exporter, _ := prometheus.New(
		prometheus.WithoutCounterSuffixes(),
		prometheus.WithoutUnits(),
	)

	// OpenTelemetry provider.
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
	)
	provider.SetAsGlobal()
	defer provider.Shutdown(ctx)

	// Counter.

	s := g.Server()
	// Root route for backward compatibility
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello from otel-go-webapp (goframe)")
	})
	// Hello World endpoint
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)

	s.SetPort(8080)
	s.Run()
}
