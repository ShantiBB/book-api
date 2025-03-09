package main

import (
	book_query "book/internal/storage/book"
	"log/slog"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"book/internal/config"
	"book/internal/lib/sl"
	"book/internal/models"
	"book/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)

	log.Debug("Debug is true")
	log.Info("Running server", slog.String("env", cfg.Env))
	log.Info(
		"HTTP Server",
		slog.String("address", cfg.HTTPServer.Address),
		slog.String("port", cfg.HTTPServer.Port),
	)

	db, err := storage.SessionDB(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize database", err)
		os.Exit(1)
	}
	log.Info("Database initialized successfully")

	book := models.Book{
		Title:       "harry potter",
		Description: "The boy why survived",
		Author:      "troll",
	}

	if err = book_query.Create(&book, db); err != nil {
		log.Error("Error book", "error", err)
		os.Exit(1)
	}
	log.Info(
		"Book created successfully",
		slog.String("id", strconv.FormatInt(book.ID, 10)),
		slog.String("title", book.Title),
	)

	defer func() {
		if err := storage.CloseDB(db); err != nil {
			log.Error("Failed to close database", err)
		}
		log.Info("Database closed successfully")
	}()
	// TODO: Add CRUD for books
}
