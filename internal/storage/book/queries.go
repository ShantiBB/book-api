package book_query

import (
	"database/sql"
	"fmt"
	"time"

	"book/internal/models"
	"book/internal/storage/queries"
)

type BookQuery interface {
	RetrieveAll()
	Retrieve()
	Create()
	Update()
	Delete()
}

func RetrieveAll(dbPath string) ([]BookQuery, error) {
	const op = "storage.sqlite.RetrieveAll"
	return nil, nil
}

func Retrieve(dbPath string) (BookQuery, error) {
	const op = "storage.sqlite.Retrieve"
	return nil, nil
}

func Create(book *models.Book, db *sql.DB) error {
	const op = "storage.sqlite.Create"

	stmt, err := db.Prepare(queries.CreateBook)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare book query - %w", op, err)
	}

	var result sql.Result
	result, err = stmt.Exec(book.Title, book.Description, book.Author)
	if err != nil {
		return fmt.Errorf("%s: failed to execute book query - %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: failed to get last insert ID - %w", op, err)
	}

	row := db.QueryRow(queries.GetCreatedAtBook, id)
	var createdAt time.Time
	if err := row.Scan(&createdAt); err != nil {
		return fmt.Errorf("%s: failed to fetch created_at - %w", op, err)
	}

	book.ID = id
	book.CreatedAt = createdAt

	return nil
}

func Update(dbPath string, pk int) (BookQuery, error) {
	const op = "storage.sqlite.Update"
	return nil, nil
}

func Delete(dbPath string, pk int) (BookQuery, error) {
	const op = "storage.sqlite.Delete"
	return nil, nil
}
