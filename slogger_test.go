package slogger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	slogotel "github.com/remychantenay/slog-otel"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func TestNewLoggerOpts(t *testing.T) {
	opts := NewLoggerOpts("testService", "testApp")
	assert.NotNil(t, opts)
	assert.Len(t, opts.addAttr, 1)
	assert.Equal(t, "system", opts.addAttr[0].Key)
	assert.Contains(t, opts.addAttr[0].String(), "testService")
	assert.Contains(t, opts.addAttr[0].String(), "testApp")

}

func TestWithDestination(t *testing.T) {
	var buf bytes.Buffer
	opt := WithDestination(&buf)
	opts := &LoggerOpts{}
	opt(opts)
	assert.Equal(t, &buf, opts.destination)
}

func TestWithJsonOutput(t *testing.T) {
	opt := WithJsonOutput()
	opts := &LoggerOpts{}
	opt(opts)
	assert.Equal(t, ModeJSON, opts.mode)
}

func TestWithAttr(t *testing.T) {
	attr := slog.String("test", "value")
	opt := WithAttr(attr)
	opts := &LoggerOpts{}
	opt(opts)
	assert.Contains(t, opts.addAttr, attr)
}

func TestNewLogger(t *testing.T) {
	opts := NewLoggerOpts("testService", "testApp")
	logger := NewLogger(opts)
	assert.NotNil(t, logger)
}

func TestWithSource(t *testing.T) {
	opt := WithSource()
	opts := &slog.HandlerOptions{}
	opt(opts)
	assert.True(t, opts.AddSource)
}

func TestWithLevel(t *testing.T) {
	opt := WithLevel("debug")
	opts := &slog.HandlerOptions{}
	opt(opts)
	assert.Equal(t, slog.LevelDebug, opts.Level)
}

func TestWithReplaceAttr(t *testing.T) {
	fn := func(groups []string, a slog.Attr) slog.Attr { return a }
	opt := WithReplaceAttr(fn)
	opts := &slog.HandlerOptions{}
	assert.Nil(t, opts.ReplaceAttr)
	opt(opts)
	assert.NotNil(t, opts.ReplaceAttr)
}

func TestLevelMapper(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, levelMapper("debug"))
	assert.Equal(t, slog.LevelInfo, levelMapper("info"))
	assert.Equal(t, slog.LevelWarn, levelMapper("warn"))
	assert.Equal(t, slog.LevelError, levelMapper("error"))
	assert.Equal(t, slog.LevelInfo, levelMapper("unknown"))
}

func TestSubLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(NewLoggerOpts("example", "api", WithDestination(&buf)))
	l := SubLogger(logger, "testComponent")
	l.Info("hello-world")
	data := buf.Bytes()
	assert.Contains(t, string(data), "testComponent")
}

func TestNewLoggerWithJsonOutput(t *testing.T) {
	t.Run("it should create the logs correctly", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithJsonOutput(), WithDestination(&buf))
		logger := NewLogger(opts)
		assert.NotNil(t, logger)

		// Log something
		logger.Info("test")

		// Check if the logged data is in JSON format
		data := buf.String()
		assert.Contains(t, data, "{")
		assert.Contains(t, data, "}")
		assert.Contains(t, data, "\"level\":\"INFO\"")
		assert.Contains(t, data, "\"msg\":\"test\"")
	})
	t.Run("it should not duplicate keys", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithJsonOutput(), WithDestination(&buf))
		logger := NewLogger(opts)
		assert.NotNil(t, logger)
		// log the same key twice
		logger.Info("test", slog.String("key", "value"), slog.String("key", "new-value"))
		data := buf.String()
		assert.Equal(t, 1, strings.Count(data, `"key":"new-value"`))
	})

}

func TestNewLoggerWithHandlerOpts(t *testing.T) {
	t.Run("it should create the logs correctly", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithDestination(&buf))
		logger := NewLogger(opts, WithSource(), WithLevel("debug"))

		assert.NotNil(t, logger)

		// Log something
		logger.Debug("test")

		// Check if the logged data contains the source and level
		data := buf.String()
		assert.Contains(t, data, "source")
		assert.Contains(t, data, "level=DEBUG")
	})
	t.Run("it should not duplicate keys", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithDestination(&buf))
		logger := NewLogger(opts)
		assert.NotNil(t, logger)
		// log the same key twice
		logger.Info("test", slog.String("key", "value"), slog.String("key", "new-value"))
		data := buf.String()
		assert.Equal(t, 1, strings.Count(data, "key=new-value"))
	})
}

func TestWithOtel(t *testing.T) {
	t.Run("it should create the logs correctly with otel enabled", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithOtel(), WithDestination(&buf))
		logger := NewLogger(opts)
		assert.NotNil(t, logger)
		ctx, span := otel.GetTracerProvider().Tracer("test").Start(t.Context(), t.Name())
		defer span.End()
		logger.InfoContext(ctx, "test")
		assert.Contains(t, buf.String(), "\"level\":\"INFO\"")
	})
}

func TestWithOtelLevel(t *testing.T) {
	t.Run("it should create the logs correctly with otel enabled", func(t *testing.T) {
		var buf bytes.Buffer
		opts := NewLoggerOpts("testService", "testApp", WithOtelOpts(slogotel.WithNoTraceEvents(true)), WithDestination(&buf))
		logger := NewLogger(opts)
		assert.NotNil(t, logger)
		ctx, span := otel.GetTracerProvider().Tracer("test").Start(t.Context(), t.Name())
		defer span.End()
		logger.InfoContext(ctx, "test")
		assert.Contains(t, buf.String(), "\"level\":\"INFO\"")

	})
}
