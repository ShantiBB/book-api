package request

type Book struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Author      string  `json:"author" validate:"required"`
	Genres      []int64 `json:"genres" validate:"required,min=1"`
}

type BookUpdate struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Genres      []int64 `json:"genres" validate:"min=1"`
}
