package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"book-api/internal/entity"
	"book-api/internal/http/schema/request"
	"book-api/internal/http/schema/response"
	"book-api/internal/http/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type BookService interface {
	CreateBook(ctx context.Context, b *entity.CreateBook) error
	GetBookByID(ctx context.Context, id int64) (*entity.Book, error)
	GetAllBooks(ctx context.Context) ([]*entity.Book, error)
	UpdateBookByID(ctx context.Context, b *entity.UpdateBook) (*entity.UpdateBook, error)
	DeleteBookByID(ctx context.Context, id int64) error
}

func (h *Handler) CreateBook() http.HandlerFunc {
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

		book := &entity.CreateBook{
			Title:       req.Title,
			Description: req.Description,
			Author:      req.Author,
			GenreIDs:    req.Genres,
		}

		if err := h.svc.CreateBook(r.Context(), book); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create book"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, book)
	}
}

func (h *Handler) GetBookByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		book, err := h.svc.GetBookByID(r.Context(), id)

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

func (h *Handler) GetBookAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.svc.GetAllBooks(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get books"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, books)
	}
}

func (h *Handler) UpdateBookByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
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

		book := &entity.UpdateBook{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			Author:      req.Author,
			GenreIDs:    req.Genres,
		}

		book, err = h.svc.UpdateBookByID(r.Context(), book)
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
		render.JSON(w, r, book)
	}
}

func (h *Handler) DeleteBookByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		if err = h.svc.DeleteBookByID(r.Context(), id); err != nil {
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
