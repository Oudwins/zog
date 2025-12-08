package logger

import (
	"io"
	"log/slog"
	"os"
)

func Init() {
	// Create log file
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Create multi writer to log to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// // Create and set default logger with no additional formatting
	// logger := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
	// 	Level: slog.LevelInfo,
	// }))

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}
