package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"book/internal/storage/queries"
)

func InitDB(storagePath string) (*sql.DB, error) {
	const op = "storage.sqlite.Init"
	var db *sql.DB

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open sqlite3 file %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	stmt, err := db.Prepare(queries.CreateBookTable)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare query - %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query - %w", op, err)
	}

	return db, nil
}
