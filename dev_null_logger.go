package slogger

import "log/slog"

// DevNullLogger returns a logger that writes to nowhere
// this is useful for testing
type DevNullLogger struct{}

// Write implements the io.Writer interface but does nothing
func (d *DevNullLogger) Write(p []byte) (n int, err error) { return }

// NewDevNullLogger returns a logger that writes to nowhere
// this is useful for testing
func NewDevNullLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&DevNullLogger{}, &slog.HandlerOptions{}))
}
