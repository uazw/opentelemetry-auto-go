package main

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/os/glog"
)

// JsonOutputsForLogger is for JSON marshaling in sequence.
type HandlerOutputJson struct {
	Timestamp  string `json:"timestamp"`                                      // Formatted time string, like "2016-01-09 12:00:00".
	TraceId    string `json:"logging.googleapis.com/trace,omitempty"`         // Trace id, only available if tracing is enabled.
	Sampled    bool   `json:"logging.googleapis.com/trace_sampled,omitempty"` // Trace id, only available if tracing is enabled.
	Level      string `json:"severity"`                                       // Formatted level string, like "DEBU", "ERRO", etc. Eg: ERRO
	CallerPath string `json:"callerPath,omitempty"`                           // The source file path and its line number that calls logging, only available if F_FILE_SHORT or F_FILE_LONG set.
	CallerFunc string `json:"callerFunc,omitempty"`                           // The source function name that calls logging, only available if F_CALLER_FN set.
	Content    string `json:"message"`                                        // Content is the main logging content, containing error stack string produced by logger.
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
	output := HandlerOutputJson{
		Timestamp:  in.TimeFormat,
		TraceId:    in.TraceId,
		Sampled:    true,
		Level:      levelMapping(in.LevelFormat),
		CallerFunc: in.CallerFunc,
		CallerPath: in.CallerPath,
		Content:    in.Content,
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
