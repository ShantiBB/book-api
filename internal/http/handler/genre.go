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

type GenreService interface {
	CreateGenre(ctx context.Context, g *entity.Genre) error
	GetGenreByID(ctx context.Context, id int64) (*entity.Genre, error)
	GetAllGenres(ctx context.Context) ([]*entity.Genre, error)
	DeleteGenreByID(ctx context.Context, id int64) error
}

func (h *Handler) CreateGenre() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.Genre

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

		genre := &entity.Genre{
			Name: req.Name,
		}

		if err := h.svc.CreateGenre(r.Context(), genre); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create genre"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, genre)
	}
}

func (h *Handler) GetGenreByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		genre, err := h.svc.GetGenreByID(r.Context(), id)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("genre not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get genre"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, genre)
	}
}

func (h *Handler) GetGenreAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genres, err := h.svc.GetAllGenres(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get genres"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, genres)
	}
}

func (h *Handler) DeleteGenreByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseID(w, r, chi.URLParam(r, "id"))
		if err != nil {
			return
		}

		if err = h.svc.DeleteGenreByID(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("book not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete genre"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
