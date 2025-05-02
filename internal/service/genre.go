package service

import (
	"context"

	"book-api/internal/entity"
)

type GenreRepository interface {
	CreateGenre(ctx context.Context, g *entity.Genre) (*entity.Genre, error)
	GetGenreByID(ctx context.Context, id int64) (*entity.Genre, error)
	GetAllGenres(ctx context.Context) ([]*entity.Genre, error)
	DeleteGenreByID(ctx context.Context, id int64) error
}

func (s *Service) CreateGenre(ctx context.Context, g *entity.Genre) error {
	const op = "genre.service.Create"

	genre, err := s.repo.CreateGenre(ctx, g)
	if err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	s.log.Info("success", "op", op, "id", genre.ID)

	return nil
}

func (s *Service) GetGenreByID(ctx context.Context, id int64) (*entity.Genre, error) {
	const op = "genre.service.GetByID"

	genre, err := s.repo.GetGenreByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get book by id", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", genre.ID)

	return genre, nil
}

func (s *Service) GetAllGenres(ctx context.Context) ([]*entity.Genre, error) {
	const op = "genre.service.GetAll"

	genres, err := s.repo.GetAllGenres(ctx)
	if err != nil {
		s.log.Error("failed to get all books", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "count", len(genres))
	return genres, nil
}

func (s *Service) DeleteGenreByID(ctx context.Context, id int64) error {
	const op = "genre.service.DeleteByID"

	if err := s.repo.DeleteGenreByID(ctx, id); err != nil {
		s.log.Error("failed to delete book", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", id)

	return nil
}
