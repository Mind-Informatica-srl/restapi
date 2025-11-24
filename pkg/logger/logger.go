// Package logger is the logging infrastructure of the project
package logger

import (
	"log/slog"
	"os"
)

type leveler struct {
	production bool
}

func (l leveler) Level() slog.Level {
	if l.production {
		return slog.LevelInfo
	}
	return slog.LevelDebug
}

var logger *slog.Logger

// Setup initialize the logging infrastructure using production mode if needed
func Setup(production bool) {
	leveler := leveler{production: production}
	opts := &slog.HandlerOptions{
		Level: leveler,
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

// Log returns the root logger
func Log() *slog.Logger {
	return logger
}
