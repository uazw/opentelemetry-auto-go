package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/prometheus/otlptranslator"

	"go.opentelemetry.io/otel/exporters/prometheus"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
)

const (
	instrument        = "github.com/this/is/a/example/basic"
	instrumentVersion = "v1.0"
)

func main() {
	var ctx = gctx.New()
	// Set custom JSON logging handler to match google structured logging format.
	glog.SetDefaultHandler(LoggingJsonHandler)

	// Set up OpenTelemetry Prometheus exporter to export metric as prometheus format.
	exporter, _ := prometheus.New(prometheus.WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes))
	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
	)

	provider.SetAsGlobal()
	defer provider.Shutdown(ctx)

	// Counter.

	s := g.Server()
	// Root route for backward compatibility
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "helworld")
		r.Response.Write("hello from otel-go-webapp (goframe)")
	})
	// Hello World endpoint
	s.BindHandler("/hello", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "helworldworld")
		r.Response.Write("hello world")
	})

	// Prometheus metrics endpoint
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)

	s.SetPort(8080)
	s.Run()
}
