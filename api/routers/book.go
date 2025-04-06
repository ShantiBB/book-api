package routers

import (
	bookQuery "book/internal/storage/book"
	"net/http"

	"github.com/gin-gonic/gin"

	"book/internal/models"
)

type Handler struct {
	bookStorage bookQuery.BookQuery
}

func NewHandler(bookStorage bookQuery.BookQuery) *Handler {
	return &Handler{bookStorage: bookStorage}
}

func (h *Handler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.bookStorage.Create(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *Handler) GetAllBooks(c *gin.Context) {
	books, err := h.bookStorage.RetrieveAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get books"})
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(c *gin.Context) {
	strID := c.Param("id")
	id, err := parseStringID(strID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	book, err := h.bookStorage.Retrieve(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *Handler) UpdateBookByID(c *gin.Context) {
	strID := c.Param("id")
	id, err := parseStringID(strID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	book, err := h.bookStorage.Retrieve(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	}

	updateBook := models.UpdateBookRequest{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
	}
	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.bookStorage.Update(&updateBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
	}

	book.Title = updateBook.Title
	book.Description = updateBook.Description
	c.JSON(http.StatusAccepted, book)
}

func (h *Handler) DeleteBookByID(c *gin.Context) {
	strID := c.Param("id")
	id, err := parseStringID(strID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	if err := h.bookStorage.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
	}

	c.JSON(http.StatusNoContent, nil)
}
