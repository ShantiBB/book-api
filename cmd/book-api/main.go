package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	bookHandler "testSQLC/internal/http/handler"
	router "testSQLC/internal/http/router/chi"
	bookRepo "testSQLC/internal/repository/postgres"
	bookService "testSQLC/internal/service"
	"testSQLC/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func main() {
	logHandler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)

	log := slog.New(logHandler)
	log.Info("init logger")

	db, err := postgres.NewPool(log)
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		err = postgres.DBClose(db, log)
		if err != nil {
			os.Exit(1)
		}

		log.Info("success close database")
	}()

	bookSvc := bookService.New(db, log, bookRepo.New(db))
	bookHdlr := bookHandler.New(db, log, bookSvc)

	r := chi.NewRouter()
	router.New(r, bookHdlr)

	log.Info("start server", "address", "localhost:8080")
	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      r,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	log.Error("server stopped")

}
