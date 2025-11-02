package main

import (
	"context"

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
	providerShutdown, _ := initMeterProvider()

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
		gauge = meter.MustUpDownCounter(
			"goframe.metric.demo.gauge",
			gmetric.MetricOption{
				Help: "This is a simple demo for UpDownCounter usage",
				Unit: "%",
			},
		)
		histogram = meter.MustHistogram(
			"goframe.metric.demo.histogram",
			gmetric.MetricOption{
				Help:    "This is a simple demo for histogram usage",
				Unit:    "ms",
				Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			},
		)
	)
	
	s := g.Server()
	// Root route for backward compatibility
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "helworld")

		counter.Add(r.Context(), 1)

		gauge.Add(r.Context(), 10) // Add adds the given value to the counter. It panics if the value is < 0
		gauge.Dec(r.Context())

		histogram.Record(1)
		histogram.Record(20)
		histogram.Record(30)
		histogram.Record(101)
		histogram.Record(2000)
		histogram.Record(9000)
		histogram.Record(20000)

		r.Response.Write("hello from otel-go-webapp (goframe)")
	})
	// Hello World endpoint
	s.BindHandler("/hello", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), "helworldworld")
		r.Response.Write("hello world")
	})
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)

	s.SetPort(8080)
	s.Run()

	defer providerShutdown(ctx)
}

func initMeterProvider() (func(context.Context) error, error) {
	exporter, _ := prometheus.New(prometheus.WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes))

	provider := otelmetric.MustProvider(
		otelmetric.WithReader(exporter),
	)
	provider.SetAsGlobal()

	return provider.Shutdown, nil
}
