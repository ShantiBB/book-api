package queries

const CreateBookTable = `
	CREATE TABLE IF NOT EXISTS book (
	id          INTEGER PRIMARY KEY,
	title       VARCHAR(20),
	description TEXT,
	created_at  INTEGER,
	author      VARCHAR(20)
)`
