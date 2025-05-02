package postgres

import (
	"context"
	"database/sql"

	"book-api/internal/entity"
)

func (r *Repository) CreateGenre(ctx context.Context, g *entity.Genre) (*entity.Genre, error) {
	query := `INSERT INTO genres (name) 
			  VALUES ($1) RETURNING id`

	err := r.DB.QueryRow(ctx, query, g.Name).Scan(&g.ID)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (r *Repository) GetGenreByID(ctx context.Context, id int64) (*entity.Genre, error) {
	query := `SELECT id, name 
			  FROM genres 
			  WHERE id = $1`

	var g entity.Genre

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&g.ID,
		&g.Name,
	)

	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *Repository) GetAllGenres(ctx context.Context) ([]*entity.Genre, error) {
	query := `SELECT id, name
			  FROM genres`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*entity.Genre

	for rows.Next() {
		var g entity.Genre
		err = rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *Repository) DeleteGenreByID(ctx context.Context, id int64) error {
	query := `DELETE FROM genres 
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
