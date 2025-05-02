package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	bookEntity "testSQLC/internal/entity"
	"testSQLC/internal/schema/request"
	"testSQLC/internal/schema/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Create(ctx context.Context, b *bookEntity.CreateBook) error
	GetByID(ctx context.Context, id int64) (*bookEntity.Book, error)
	GetAll(ctx context.Context) ([]*bookEntity.Book, error)
	UpdateByID(ctx context.Context, b *bookEntity.UpdateBook) (*bookEntity.UpdateBook, error)
	DeleteByID(ctx context.Context, id int64) error
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.Book

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		book := &bookEntity.CreateBook{
			Title:       req.Title,
			Description: req.Description,
			Author:      req.Author,
			GenreIDs:    req.Genres,
		}

		if err := h.svc.Create(r.Context(), book); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create book"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, response.Book{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Author:      book.Author,
			Genres:      book.GenreIDs,
		})
	}
}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")

		id, err := strconv.Atoi(strID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		book, err := h.svc.GetByID(r.Context(), int64(id))

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("book not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get book"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, book)
	}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.svc.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get books"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, books)
	}
}

func (h *Handler) UpdateByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")

		id, err := strconv.Atoi(strID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		var req request.BookUpdate
		if err = render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to render"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		book := &bookEntity.UpdateBook{
			ID:          int64(id),
			Title:       req.Title,
			Description: req.Description,
			Author:      req.Author,
			GenreIDs:    req.Genres,
		}

		book, err = h.svc.UpdateByID(r.Context(), book)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("book not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update book"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Book{
			ID:          int64(id),
			Title:       book.Title,
			Description: book.Description,
			Author:      book.Author,
			Genres:      book.GenreIDs,
		})
	}
}

func (h *Handler) DeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")

		id, err := strconv.Atoi(strID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		if err = h.svc.DeleteByID(r.Context(), int64(id)); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("book not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete book"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
