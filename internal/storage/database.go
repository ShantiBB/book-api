package storage

import (
	"database/sql"
	"fmt"

	"book/internal/storage/queries"
)

const op = "storage.sqlite.Init"

func CheckDB(db *sql.DB, err error) error {
	if err != nil {
		return fmt.Errorf("%s: failed to open sqlite3 file %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return nil
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

func SessionDB(storagePath string) (*sql.DB, error) {
	var db *sql.DB

	db, err := sql.Open("sqlite", storagePath)
	ok := CheckDB(db, err)
	if ok != nil {
		return nil, ok
	}

	stmt, err := db.Prepare(queries.CreateBookTable)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare book table query - %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed book table - %w", op, err)
	}

	return db, nil
}
