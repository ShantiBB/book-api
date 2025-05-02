package service

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB   *pgxpool.Pool
	Log  *slog.Logger
	Repo Repository
}

func New(db *pgxpool.Pool, log *slog.Logger, repo Repository) *Service {
	return &Service{
		DB:   db,
		Log:  log,
		Repo: repo,
	}
}
