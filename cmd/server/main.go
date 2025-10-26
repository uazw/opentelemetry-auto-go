package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmetric"
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
	glog.SetDefaultHandler(LoggingJsonHandler)

	exporter, _ := prometheus.New(prometheus.WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes))

	var (
		meter = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
			Instrument:        instrument,
			InstrumentVersion: instrumentVersion,
		})
		counter = meter.MustCounter(
			"goframe.metric.demo.counter",
			gmetric.MetricOption{
				Help: "This is a simple demo for Counter usage",
				Unit: "bytes",
			},
		)
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
		g.Log().Info(r.Context(), "helworld")
		r.Response.Write("hello from otel-go-webapp (goframe)")
	})
	// Hello World endpoint
	s.BindHandler("/hello", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "helworldworld")
		r.Response.Write("hello world")
	})
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)
	counter.Add(ctx, 1)

	s.SetPort(8080)
	s.Run()
}
