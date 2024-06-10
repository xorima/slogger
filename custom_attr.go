package slogger

import "log/slog"

func ErrorAttr(err error) slog.Attr {
	if err != nil {
		return slog.String("error", err.Error())
	}
	return slog.String("error", "")
}
