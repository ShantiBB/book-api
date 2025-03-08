package sl

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	levelDebug := slog.HandlerOptions{Level: slog.LevelDebug}
	levelInfo := slog.HandlerOptions{Level: slog.LevelInfo}

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &levelDebug),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &levelDebug),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &levelInfo),
		)
	}

	return log
}
