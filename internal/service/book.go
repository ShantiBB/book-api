package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	"testSQLC/internal/entity"
	"testSQLC/internal/utils"
)

type Repository interface {
	CreateBook(ctx context.Context, tx pgx.Tx, b *entity.CreateBook) error
	GetBookByID(ctx context.Context, id int64) (*entity.Book, error)
	GetAllBooks(ctx context.Context) ([]*entity.Book, error)
	UpdateBookByID(ctx context.Context, tx pgx.Tx, b *entity.UpdateBook) (*entity.UpdateBook, error)
	DeleteBookByID(ctx context.Context, id int64) error
	AddGenres(ctx context.Context, tx pgx.Tx, id int64, genres []int64) error
	GetAllBooksGenres(ctx context.Context) (map[int64][]entity.Genre, error)
	GetGenresByBookID(ctx context.Context, id int64) ([]entity.Genre, error)
	UpdateBookGenres(ctx context.Context, tx pgx.Tx, id int64, genres []int64) ([]int64, error)
}

func (s *Service) Create(ctx context.Context, b *entity.CreateBook) error {
	const op = "book.service.CreateBook"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Log.Error("failed to begin transaction", "op", op, "error", err)
		return err
	}

	defer func() {
		utils.CloseTransaction(ctx, tx, err)
	}()

	if err = s.Repo.CreateBook(ctx, tx, b); err != nil {
		s.Log.Error("failed", "op", op, "error", err)
		return err
	}

	if err = s.Repo.AddGenres(ctx, tx, b.ID, b.GenreIDs); err != nil {
		s.Log.Error("failed to add genres", "op", op, "error", err)
		return err
	}

	s.Log.Debug("success", "op", op, "id", b.ID)

	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*entity.Book, error) {
	const op = "book.service.GetBookByID"

	book, err := s.Repo.GetBookByID(ctx, id)
	if err != nil {
		s.Log.Error("failed to get book by id", "op", op, "error", err)
		return nil, err
	}

	book.Genres, err = s.Repo.GetGenresByBookID(ctx, id)
	if err != nil {
		s.Log.Error("failed to get genres by id", "op", op, "error", err)
		return nil, err
	}

	s.Log.Debug("success", "op", op, "id", book.ID)

	return book, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*entity.Book, error) {
	const op = "book.service.GetAllBooks"

	books, err := s.Repo.GetAllBooks(ctx)
	if err != nil {
		s.Log.Error("failed to get all books", "op", op, "error", err)
		return nil, err
	}

	genresMap, err := s.Repo.GetAllBooksGenres(ctx)
	if err != nil {
		s.Log.Error("failed to get all genres", "op", op, "error", err)
		return nil, err
	}

	for _, book := range books {
		book.Genres = genresMap[book.ID]
	}

	s.Log.Debug("success", "op", op, "count", len(books))
	return books, nil
}

func (s *Service) UpdateByID(ctx context.Context, b *entity.UpdateBook) (*entity.UpdateBook, error) {
	const op = "book.service.Update"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Log.Error("failed to begin transaction", "op", op, "error", err)
		return nil, err
	}

	defer func() {
		utils.CloseTransaction(ctx, tx, err)
	}()

	updatedBook, err := s.Repo.UpdateBookByID(ctx, tx, b)
	if err != nil {
		s.Log.Error("failed to update book", "op", op, "error", err)
		return nil, err
	}

	updatedBook.GenreIDs, err = s.Repo.UpdateBookGenres(ctx, tx, b.ID, b.GenreIDs)
	if err != nil {
		s.Log.Error("failed to update book genres", "op", op, "error", err)
		return nil, err
	}

	s.Log.Debug("success", "op", op, "id", b.ID)

	return updatedBook, nil
}

func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	const op = "book.service.DeleteBookByID"

	if err := s.Repo.DeleteBookByID(ctx, id); err != nil {
		s.Log.Error("failed to delete book", "op", op, "error", err)
		return err
	}

	s.Log.Debug("success", "op", op, "id", id)

	return nil
}
