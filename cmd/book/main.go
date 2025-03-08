package main

import (
	"database/sql"
	"log/slog"

	"book/internal/config"
	"book/internal/lib/sl"
	"book/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)

	log.Info("Running server", slog.String("env", cfg.Env))
	log.Debug("Debug is true")
	log.Info("HTTP Server", slog.String("address", cfg.HTTPServer.Address))
	log.Debug("HTTP Server", slog.String("port", cfg.HTTPServer.Port))

	db, err := storage.InitDB(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize database: %v", err)
	}

	log.Info("Database initialized successfully")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error("Failed to close database: %v", err)
		}
	}(db)
}
