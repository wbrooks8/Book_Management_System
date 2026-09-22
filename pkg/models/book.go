package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

// Repository contains database operations for books.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a book repository backed by db.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Migrate creates or updates the books table.
func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&Book{})
}

func (r *Repository) CreateBook(book *Book) error {
	return r.db.Create(book).Error
}

func (r *Repository) GetAllBooks() ([]Book, error) {
	var books []Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *Repository) GetBookByID(id int64) (*Book, error) {
	var book Book
	err := r.db.First(&book, id).Error
	return &book, err
}

func (r *Repository) UpdateBook(book *Book) error {
	return r.db.Save(book).Error
}

func (r *Repository) DeleteBook(id int64) (*Book, error) {
	var book Book
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, r.db.Delete(&book).Error
}
