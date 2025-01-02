package logging

import (
	"log/slog"
	"os"
	"strings"
)

var (
	Log     *slog.Logger
	Handler *slog.JSONHandler
	Level   = &slog.LevelVar{}
)

func init() {
	Handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: Level,
	})

	Log = slog.New(Handler)
}

func SetLogLevel(l string) slog.Level {
	switch strings.ToLower(l) {
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
