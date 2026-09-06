package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *slog.Logger {
	logFilePath := filepath.Join("logs", "app.log")

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	fileWriter := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	return slog.New(slog.NewJSONHandler(multiWriter, opts))
}
