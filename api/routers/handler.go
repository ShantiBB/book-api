package routers

import (
	bookQuery "book/internal/storage/book"
)

type Handler struct {
	bookStorage bookQuery.BookQuery
}

func NewHandler(bookStorage bookQuery.BookQuery) *Handler {
	return &Handler{bookStorage: bookStorage}
}
