package chi

import (
	"github.com/go-chi/chi/v5"

	"book-api/internal/http/handler"
)

func genreRouter(h *handler.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.CreateGenre())
		r.Get("/{id}", h.GetGenreByID())
		r.Get("/", h.GetGenreAll())
		r.Delete("/{id}", h.DeleteGenreByID())
	}
}
