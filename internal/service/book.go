package service

import (
	"context"

	"book-api/internal/entity"
	"book-api/internal/utils"

	"github.com/jackc/pgx/v5"
)

type BookRepository interface {
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

func (s *Service) CreateBook(ctx context.Context, b *entity.CreateBook) error {
	const op = "book.service.Create"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", "op", op, "error", err)
		return err
	}

	defer func() {
		utils.CloseTransaction(ctx, tx, err)
	}()

	if err = s.repo.CreateBook(ctx, tx, b); err != nil {
		s.log.Error("failed", "op", op, "error", err)
		return err
	}

	if err = s.repo.AddGenres(ctx, tx, b.ID, b.GenreIDs); err != nil {
		s.log.Error("failed to add genres", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", b.ID)

	return nil
}

func (s *Service) GetBookByID(ctx context.Context, id int64) (*entity.Book, error) {
	const op = "book.service.GetByID"

	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get book by id", "op", op, "error", err)
		return nil, err
	}

	book.Genres, err = s.repo.GetGenresByBookID(ctx, id)
	if err != nil {
		s.log.Error("failed to get genres by id", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", book.ID)

	return book, nil
}

func (s *Service) GetAllBooks(ctx context.Context) ([]*entity.Book, error) {
	const op = "book.service.GetAll"

	books, err := s.repo.GetAllBooks(ctx)
	if err != nil {
		s.log.Error("failed to get all books", "op", op, "error", err)
		return nil, err
	}

	genresMap, err := s.repo.GetAllBooksGenres(ctx)
	if err != nil {
		s.log.Error("failed to get all genres", "op", op, "error", err)
		return nil, err
	}

	for _, book := range books {
		book.Genres = genresMap[book.ID]
	}

	s.log.Debug("success", "op", op, "count", len(books))
	return books, nil
}

func (s *Service) UpdateBookByID(ctx context.Context, b *entity.UpdateBook) (*entity.UpdateBook, error) {
	const op = "book.service.Update"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", "op", op, "error", err)
		return nil, err
	}

	defer func() {
		utils.CloseTransaction(ctx, tx, err)
	}()

	updatedBook, err := s.repo.UpdateBookByID(ctx, tx, b)
	if err != nil {
		s.log.Error("failed to update book", "op", op, "error", err)
		return nil, err
	}

	updatedBook.GenreIDs, err = s.repo.UpdateBookGenres(ctx, tx, b.ID, b.GenreIDs)
	if err != nil {
		s.log.Error("failed to update book genres", "op", op, "error", err)
		return nil, err
	}

	s.log.Debug("success", "op", op, "id", b.ID)

	return updatedBook, nil
}

func (s *Service) DeleteBookByID(ctx context.Context, id int64) error {
	const op = "book.service.DeleteByID"

	if err := s.repo.DeleteBookByID(ctx, id); err != nil {
		s.log.Error("failed to delete book", "op", op, "error", err)
		return err
	}

	s.log.Debug("success", "op", op, "id", id)

	return nil
}
