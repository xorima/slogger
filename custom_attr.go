package slogger

import "log/slog"

// ErrorAttr writes the error (if it exists!) to the key of `error`
// this ensures a consistent approach to error log keys across
// an application
func ErrorAttr(err error) slog.Attr {
	if err != nil {
		return slog.String("error", err.Error())
	}
	return slog.String("error", "")
}
