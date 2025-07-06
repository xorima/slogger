package slogger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	slogotel "github.com/remychantenay/slog-otel"
	slogdedup "github.com/veqryn/slog-dedup"
)

const (
	ModeText = "text"
	ModeJSON = "json"
)

type LoggerOpts struct {
	addAttr     []slog.Attr
	destination io.Writer
	mode        string
	otelMode    bool
	otelOpts    []slogotel.OtelHandlerOpt
}

func NewLoggerOpts(serviceName, applicationName string, opts ...func(o *LoggerOpts)) *LoggerOpts {
	o := &LoggerOpts{
		destination: os.Stdout,
		mode:        ModeText,
	}
	o.addAttr = append(o.addAttr, slog.Group("system",
		slog.String("service", serviceName),
		slog.String("applicationName", applicationName)))

	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithDestination(destination io.Writer) func(o *LoggerOpts) {
	return func(o *LoggerOpts) {
		o.destination = destination
	}
}

func WithJsonOutput() func(o *LoggerOpts) {
	return func(o *LoggerOpts) {
		o.mode = ModeJSON
	}
}

func WithAttr(attr slog.Attr) func(o *LoggerOpts) {
	return func(o *LoggerOpts) {
		o.addAttr = append(o.addAttr, attr)
	}
}

// WithOtel enables the otel logging system on this logger
// which uses github.com/remychantenay/slog-otel to add logs
// to spans and traces
func WithOtel() func(o *LoggerOpts) {
	return func(o *LoggerOpts) {
		o.otelMode = true
	}
}

// WithOtelOpts enables the user to setup the opts for the otel
// handler within the logger following github.com/remychantenay/slog-otel
// it also enables otel mode at the same time.
func WithOtelOpts(opts slogotel.OtelHandlerOpt) func(o *LoggerOpts) {
	return func(o *LoggerOpts) {
		WithOtel()(o) // As we are using otel opts, otel should be enabled
		if o.otelOpts == nil {
			o.otelOpts = make([]slogotel.OtelHandlerOpt, 0)
		}
		o.otelOpts = append(o.otelOpts, opts)
	}
}

// NewLogger returns a new slog logger
func NewLogger(loggerOpts *LoggerOpts, handlerOpts ...func(o *slog.HandlerOptions)) *slog.Logger {
	hOpts := &slog.HandlerOptions{}

	for _, opt := range handlerOpts {
		opt(hOpts)
	}
	log := slogdedup.NewOverwriteHandler(slog.NewTextHandler(loggerOpts.destination, hOpts), nil)

	if strings.ToLower(loggerOpts.mode) == ModeJSON {
		log = slogdedup.NewOverwriteHandler(slog.NewJSONHandler(loggerOpts.destination, hOpts), nil)
	}

	if loggerOpts.otelMode {
		return slog.New(slogotel.New(log, loggerOpts.otelOpts...))
	}
	return slog.New(log)
}

func WithSource() func(o *slog.HandlerOptions) {
	return func(o *slog.HandlerOptions) {
		o.AddSource = true
	}
}

func WithLevel(level string) func(o *slog.HandlerOptions) {
	return func(o *slog.HandlerOptions) {
		o.Level = levelMapper(level)
	}
}

func WithReplaceAttr(fn func(groups []string, a slog.Attr) slog.Attr) func(o *slog.HandlerOptions) {
	return func(o *slog.HandlerOptions) {
		o.ReplaceAttr = fn
	}
}

func levelMapper(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func SubLogger(logger *slog.Logger, componentName string) *slog.Logger {
	return logger.With(slog.Group("component", slog.String("name", componentName)))
}
