package book_query

const CreateBookTable = `
	CREATE TABLE IF NOT EXISTS book (
	id          INTEGER PRIMARY KEY,
	title       VARCHAR(20),
	description TEXT,
	created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
	author      VARCHAR(20)
)`

const CreateBook = `
	INSERT INTO book (title, description, author) 
	VALUES(?, ?, ?)
`

const GetCreatedAtBook = `SELECT created_at FROM book WHERE id = ?`

const GetBooks = `SELECT * FROM book`

const GetBookByID = `
	SELECT id, title, description, author, created_at 
	FROM book 
	WHERE id = ?
`

const UpdateBookByID = `UPDATE book SET title = ?, description = ? WHERE id = ?`

const DeleteBook = `DELETE FROM book WHERE id = ?`
