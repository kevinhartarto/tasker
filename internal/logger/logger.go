package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	Logger = slog.New(slog.NewJSONHandler(logFile, handlerOpts))
	slog.SetDefault(Logger)
	Logger.Info("Logger initialized")
}

func GetLogger() *slog.Logger {
	return Logger
}
