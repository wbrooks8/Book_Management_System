package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wbrook8/book_management_system/pkg/models"
	"gorm.io/gorm"
)

type Handler struct {
	books *models.Repository
}

func NewHandler(books *models.Repository) *Handler {
	return &Handler{books: books}
}

func (h *Handler) GetBooks(context *gin.Context) {
	books, err := h.books.GetAllBooks()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve books."})
		return
	}
	context.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(context *gin.Context) {
	bookID, err := parseBookID(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID."})
		return
	}
	book, err := h.books.GetBookByID(bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve book."})
		return
	}
	context.JSON(http.StatusOK, book)
}

func (h *Handler) CreateBook(context *gin.Context) {
	book := &models.Book{}
	if err := context.ShouldBindJSON(book); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}
	if err := h.books.CreateBook(book); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create book."})
		return
	}
	context.JSON(http.StatusCreated, book)
}

func (h *Handler) DeleteBook(context *gin.Context) {
	bookID, err := parseBookID(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID."})
		return
	}
	_, err = h.books.DeleteBook(bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete book."})
		return
	}
	context.Status(http.StatusNoContent)
}

func (h *Handler) UpdateBook(context *gin.Context) {
	updateBook := &models.Book{}
	bookID, err := parseBookID(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid book ID."})
		return
	}
	if err := context.ShouldBindJSON(updateBook); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}
	book, err := h.books.GetBookByID(bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve book."})
		return
	}
	if updateBook.Name != "" {
		book.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		book.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		book.Publication = updateBook.Publication
	}
	if err := h.books.UpdateBook(book); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update book."})
		return
	}
	context.JSON(http.StatusOK, book)
}

func parseBookID(context *gin.Context) (int64, error) {
	return strconv.ParseInt(context.Param("bookId"), 10, 64)
}
