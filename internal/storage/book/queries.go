package book_query

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"book/internal/models"
	"book/internal/storage/queries"
)

type BookQuery interface {
	Create(book *models.Book) error
	Retrieve(id int64) (string, error)
	RetrieveAll() ([]models.Book, error)
	Update(book *models.Book) error
	Delete(id int64) error
}

type SQLiteBookStorage struct {
	DB *sql.DB
}

func (s *SQLiteBookStorage) Create(book *models.Book) error {
	const op = "storage.sqlite.Create"

	res, err := s.DB.Exec(queries.CreateBook, book.Title, book.Description, book.Author)
	if err != nil {
		return fmt.Errorf("%s: failed to execute book query - %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: failed to get last insert ID - %w", op, err)
	}

	var createdAt time.Time
	if err := s.DB.QueryRow(queries.GetCreatedAtBook, id).Scan(&createdAt); err != nil {
		return fmt.Errorf("%s: failed to get created time - %w", op, err)
	}

	book.ID = id
	book.CreatedAt = createdAt
	return nil
}

func (s *SQLiteBookStorage) RetrieveAll() ([]models.Book, error) {
	const op = "storage.sqlite.RetrieveAll"

	rows, err := s.DB.Query(queries.GetBooks)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute books query - %s", op, err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.CreatedAt, &b.Author)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan book row - %s", op, err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows books iteration error - %s", op, err)
	}

	return books, nil
}

func (s *SQLiteBookStorage) Retrieve(id int64) (string, error) {
	const op = "storage.sqlite.Retrieve"
	var res string

	err := s.DB.QueryRow(queries.GetBookByID, id).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: book not found - %s", op, err)
	}
	if err != nil {
		return "", fmt.Errorf("%s: failed to fetch book %w", op, err)
	}

	return res, nil
}

func (s *SQLiteBookStorage) Update(b *models.Book) error {
	const op = "storage.sqlite.Update"

	res, err := s.DB.Exec(queries.UpdateBookByID, &b.Title, &b.Description, &b.ID)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query - %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected - %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows updated", op)
	}

	return nil
}

func (s *SQLiteBookStorage) Delete(id int64) error {
	const op = "storage.sqlite.Delete"

	res, err := s.DB.Exec(queries.DeleteBook, id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute book query - %s", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected - %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows deleted, book with id %d not found", op, id)
	}

	return nil
}
