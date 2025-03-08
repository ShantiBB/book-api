package models

type Book struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   int    `json:"year"`
	Author      string `json:"author"`
}
