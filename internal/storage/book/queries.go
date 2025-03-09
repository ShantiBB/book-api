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
	RetrieveAll()
	Retrieve()
	Create()
	Update()
	Delete()
}

func RetrieveAll(db *sql.DB) ([]models.Book, error) {
	const op = "storage.sqlite.RetrieveAll"

	stmt, err := db.Prepare(queries.GetBooks)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare books query - %s", op, err)
	}

	rows, err := stmt.Query()
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

func Retrieve(id int64, db *sql.DB) (string, error) {
	const op = "storage.sqlite.Retrieve"

	stmt, err := db.Prepare(queries.GetBookByID)
	if err != nil {
		return "", fmt.Errorf("%s: failed to prepare book query - %s", op, err)
	}

	var res string
	err = stmt.QueryRow(id).Scan(&res)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: book is not  found - %s", op, err)
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement %w", op, err)
	}

	return res, nil
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
