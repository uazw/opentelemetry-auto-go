package main

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/os/glog"
	"go.opentelemetry.io/otel/trace"
)

type HandlerOutputJson struct {
	Timestamp string `json:"timestamp"`                                      // Formatted time string, like "2016-01-09 12:00:00".
	TraceId   string `json:"logging.googleapis.com/trace,omitempty"`         // Trace id, only available if tracing is enabled.
	SpanId    string `json:"logging.googleapis.com/spanId,omitempty"`        // Trace id, only available if tracing is enabled.
	Sampled   bool   `json:"logging.googleapis.com/trace_sampled,omitempty"` // Trace id, only available if tracing is enabled.
	Level     string `json:"severity"`                                       // Formatted level string, like "DEBU", "ERRO", etc. Eg: ERRO
	Content   string `json:"message"`                                        // Content is the main logging content, containing error stack string produced by logger.
}

func levelMapping(original string) string {
	switch original {
	case "DEBU":
		return "DEBUG"
	case "ERRO":
		return "ERROR"
	case "WARN":
		return "WARNING"
	default:
		return original
	}
}

// LoggingJsonHandler is a example handler for logging JSON format content.
var LoggingJsonHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	sampled := true
	if in.TraceId == "" {
		sampled = false
	}

	SpanId := ""
	if trace.SpanFromContext(ctx).SpanContext().SpanID().IsValid() {
		SpanId = trace.SpanFromContext(ctx).SpanContext().SpanID().String()

	}

	output := HandlerOutputJson{
		Timestamp: in.TimeFormat,
		TraceId:   in.TraceId,
		SpanId:    SpanId,
		Sampled:   sampled,
		Level:     levelMapping(in.LevelFormat),
		Content:   in.Content,
	}
	if len(in.Values) > 0 {
		if output.Content != "" {
			output.Content += " "
		}
		output.Content += in.ValuesContent()
	}
	// Output json content.
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.Write([]byte("\n"))
	in.Next(ctx)
}
