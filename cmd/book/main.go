package main

import (
	"log/slog"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"book/internal/config"
	"book/internal/lib/sl"
	"book/internal/models"
	"book/internal/storage"
	bookQuery "book/internal/storage/book"
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
	log.Debug("Database initialized successfully")

	book := models.Book{
		Title:       "harry potter",
		Description: "The boy why survived",
		Author:      "Troll",
	}

	if err = bookQuery.Create(&book, db); err != nil {
		log.Error("Error book", err)
		os.Exit(1)
	}
	log.Debug(
		"Book created successfully",
		slog.String("id", strconv.FormatInt(book.ID, 10)),
		slog.String("title", book.Title),
	)

	var bookByID string
	bookByID, err = bookQuery.Retrieve(book.ID, db)
	if err != nil {
		log.Error("Error retrieving book", err)
		os.Exit(1)
	}
	log.Debug("Book retrieved successfully", slog.String("title", bookByID))

	books, err := bookQuery.RetrieveAll(db)
	if err != nil {
		log.Error("Error retrieving books", err)
		os.Exit(1)
	}

	for _, b := range books {
		log.Debug(
			"Book retrieved successfully",
			slog.String("id", strconv.FormatInt(b.ID, 10)),
			slog.String("title", b.Title),
			slog.String("author", b.Author),
		)
	}

	defer func() {
		if err := storage.CloseDB(db); err != nil {
			log.Error("Failed to close database", err)
		}
		log.Debug("Database closed successfully")
	}()
	// TODO: Add CRUD for books
}
