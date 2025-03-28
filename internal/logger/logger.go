package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	logger = slog.New(slog.NewJSONHandler(logFile, handlerOpts))
}

func GetLogger() *slog.Logger {
	return logger
}
