package postgres

import (
	"context"
	"database/sql"

	"testSQLC/internal/entity"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateBook(ctx context.Context, tx pgx.Tx, b *entity.CreateBook) error {
	query := `INSERT INTO books (title, description, author) 
			  VALUES ($1, $2, $3) RETURNING id`

	err := tx.QueryRow(ctx, query, b.Title, b.Description, b.Author).Scan(&b.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBookByID(ctx context.Context, id int64) (*entity.Book, error) {
	query := `SELECT id, title, description, author, created_at 
			  FROM books 
			  WHERE id = $1`

	var b entity.Book

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.Author,
		&b.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *Repository) GetAllBooks(ctx context.Context) ([]*entity.Book, error) {
	query := `SELECT id, title, description, author, created_at 
			  FROM books`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []*entity.Book

	for rows.Next() {
		var b entity.Book
		if err = rows.Scan(&b.ID, &b.Title, &b.Description, &b.Author, &b.CreatedAt); err != nil {
			return nil, err
		}

		books = append(books, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *Repository) UpdateBookByID(
	ctx context.Context,
	tx pgx.Tx,
	b *entity.UpdateBook) (*entity.UpdateBook, error,
) {
	query := `UPDATE books 
			  SET title = $1, description = $2, author = $3 
			  WHERE id = $4`

	res, err := tx.Exec(ctx, query, b.Title, b.Description, b.Author, b.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return b, nil
}

func (r *Repository) DeleteBookByID(ctx context.Context, id int64) error {
	query := `DELETE FROM books 
       		  WHERE id = $1`

	res, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
