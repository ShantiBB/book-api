package book_query

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"book/internal/models"
)

type BookQuery interface {
	Create(book *models.Book) error
	Retrieve(id *int) (*models.Book, error)
	RetrieveAll() ([]models.Book, error)
	Update(book *models.UpdateBookRequest) error
	Delete(id *int) error
}

type SQLiteBookStorage struct {
	DB *gorm.DB
}

func (s *SQLiteBookStorage) Create(book *models.Book) error {
	const op = "storage.sqlite.Create"

	if err := s.DB.Create(book).Error; err != nil {
		return fmt.Errorf("%s: failed to create book - %w", op, err)
	}

	return nil
}

func (s *SQLiteBookStorage) RetrieveAll() ([]models.Book, error) {
	const op = "storage.sqlite.RetrieveAll"
	var books []models.Book

	if err := s.DB.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("%s: failed to retrieve books - %w", op, err)
	}
	fmt.Println(books)
	return books, nil
}

func (s *SQLiteBookStorage) Retrieve(id *int) (*models.Book, error) {
	const op = "storage.sqlite.Retrieve"
	var book models.Book

	if err := s.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: book not found", op)
		}
		return nil, fmt.Errorf("%s: failed to retrieve book - %w", op, err)
	}

	return &book, nil
}

func (s *SQLiteBookStorage) Update(b *models.UpdateBookRequest) error {
	const op = "storage.sqlite.Update"

	result := s.DB.Model(&models.Book{}).Where("id = ?", b.ID).Updates(models.Book{
		Title:       b.Title,
		Description: b.Description,
	})
	if result.Error != nil {
		return fmt.Errorf("%s: failed to update book - %w", op, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: no rows updated", op)
	}

	return nil
}

func (s *SQLiteBookStorage) Delete(id *int) error {
	const op = "storage.sqlite.Delete"

	result := s.DB.Delete(&models.Book{}, *id)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to delete book - %w", op, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: no rows deleted", op)
	}

	return nil
}
