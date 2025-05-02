package handler

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db  *pgxpool.Pool
	log *slog.Logger
	svc Service
}

func New(db *pgxpool.Pool, log *slog.Logger, svc Service) *Handler {
	return &Handler{
		db:  db,
		log: log,
		svc: svc,
	}
}
