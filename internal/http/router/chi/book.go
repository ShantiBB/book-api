package chi

import (
	"github.com/go-chi/chi/v5"

	"book-api/internal/http/handler"
)

func bookRouter(h *handler.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.CreateBook())
		r.Get("/{id}", h.GetBookByID())
		r.Get("/", h.GetBookAll())
		r.Put("/{id}", h.UpdateBookByID())
		r.Delete("/{id}", h.DeleteBookByID())
	}
}
