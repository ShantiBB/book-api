package entity

import "time"

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Genres      []Genre   `json:"genres"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateBook struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	GenreIDs    []int64 `json:"genres"`
}

type UpdateBook struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	GenreIDs    []int64 `json:"genres"`
}
