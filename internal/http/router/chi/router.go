package chi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"testSQLC/internal/http/handler"
)

func New(r chi.Router, h *handler.Handler) {
	r.Use(middleware.CleanPath)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)

	r.Route("/books", func(r chi.Router) {
		r.Post("/", h.Create())
		r.Get("/", h.GetAll())
		r.Get("/{id}", h.GetByID())
		r.Put("/{id}", h.UpdateByID())
		r.Delete("/{id}", h.DeleteByID())
	})
}
