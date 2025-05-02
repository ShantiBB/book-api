package chi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"book-api/internal/http/handler"
)

func New(r chi.Router, h *handler.Handler) {
	r.Use(middleware.CleanPath)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)

	r.Route("/books", bookRouter(h))
	r.Route("/genres", genreRouter(h))
}
