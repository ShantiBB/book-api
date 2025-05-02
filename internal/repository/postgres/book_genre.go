package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"testSQLC/internal/entity"
)

func (r *Repository) AddGenres(ctx context.Context, tx pgx.Tx, id int64, genres []int64) error {
	genreQuery := `INSERT INTO books_genres (book_id, genre_id) VALUES ($1, $2)`

	var err error
	for _, genreID := range genres {
		_, err = tx.Exec(ctx, genreQuery, id, genreID)
		if err != nil {
			return fmt.Errorf("failed to insert genre for book %d: %w", id, err)
		}
	}

	return nil
}

func (r *Repository) GetAllBooksGenres(ctx context.Context) (map[int64][]entity.Genre, error) {
	query := `
		SELECT bg.book_id, g.id, g.name
		FROM books_genres bg
		JOIN genres g ON g.id = bg.genre_id
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var (
		genre  entity.Genre
		bookID int64
	)
	genresMap := make(map[int64][]entity.Genre)

	for rows.Next() {
		if err = rows.Scan(&bookID, &genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genresMap[bookID] = append(genresMap[bookID], genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genresMap, nil
}

func (r *Repository) GetGenresByBookID(ctx context.Context, id int64) ([]entity.Genre, error) {
	query := `
		SELECT g.id, g.name
		FROM books_genres bg
		JOIN genres g ON bg.genre_id = g.id
		WHERE bg.book_id = $1
	`

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	var (
		genre  entity.Genre
		genres []entity.Genre
	)
	for rows.Next() {
		if err = rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *Repository) UpdateBookGenres(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
	genres []int64,
) ([]int64, error) {
	deleteQuery := `DELETE FROM books_genres
					WHERE book_id = $1 AND genre_id NOT IN (
						SELECT unnest($2::bigint[])
					)`
	insertQuery := `INSERT INTO books_genres (book_id, genre_id)
					SELECT $1, g
					FROM unnest($2::bigint[]) AS g
					LEFT JOIN books_genres bg 
						ON bg.book_id = $1 AND bg.genre_id = g
					WHERE bg.genre_id IS NULL`

	var err error
	if _, err = tx.Exec(ctx, deleteQuery, id, genres); err != nil {
		return nil, err
	}

	if _, err = tx.Exec(ctx, insertQuery, id, genres); err != nil {
		return nil, err
	}

	return genres, nil
}
