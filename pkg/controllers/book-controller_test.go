package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wbrook8/book_management_system/pkg/controllers"
	"github.com/wbrook8/book_management_system/pkg/models"
	"github.com/wbrook8/book_management_system/pkg/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestServer(t *testing.T) (*gin.Engine, *models.Repository) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	books := models.NewRepository(db)
	if err := books.Migrate(); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	server := gin.New()
	routes.RegisterRoutes(server, controllers.NewHandler(books))
	return server, books
}

func TestCreateAndGetBook(t *testing.T) {
	server, _ := newTestServer(t)

	createResponse := performRequest(server, http.MethodPost, "/book/", `{"name":"Dune","author":"Frank Herbert","publication":"1965"}`)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	var created models.Book
	decodeResponse(t, createResponse, &created)
	if created.Name != "Dune" {
		t.Fatalf("created name = %q, want %q", created.Name, "Dune")
	}

	getResponse := performRequest(server, http.MethodGet, fmt.Sprintf("/book/%d", created.ID), "")
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResponse.Code, http.StatusOK)
	}
}

func TestUpdateBook(t *testing.T) {
	server, books := newTestServer(t)
	book := &models.Book{Name: "Old title", Author: "Author", Publication: "2020"}
	if err := books.CreateBook(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	response := performRequest(server, http.MethodPut, fmt.Sprintf("/book/%d", book.ID), `{"name":"New title"}`)
	if response.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", response.Code, http.StatusOK)
	}

	var updated models.Book
	decodeResponse(t, response, &updated)
	if updated.Name != "New title" {
		t.Fatalf("updated name = %q, want %q", updated.Name, "New title")
	}
	if updated.Author != "Author" {
		t.Fatalf("updated author = %q, want original author", updated.Author)
	}
}

func TestInvalidAndMissingBooks(t *testing.T) {
	server, _ := newTestServer(t)

	invalidIDResponse := performRequest(server, http.MethodGet, "/book/not-a-number", "")
	if invalidIDResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid ID status = %d, want %d", invalidIDResponse.Code, http.StatusBadRequest)
	}

	missingResponse := performRequest(server, http.MethodGet, "/book/999", "")
	if missingResponse.Code != http.StatusNotFound {
		t.Fatalf("missing book status = %d, want %d", missingResponse.Code, http.StatusNotFound)
	}

	malformedResponse := performRequest(server, http.MethodPost, "/book/", `{invalid}`)
	if malformedResponse.Code != http.StatusBadRequest {
		t.Fatalf("malformed JSON status = %d, want %d", malformedResponse.Code, http.StatusBadRequest)
	}
}

func TestDeleteBook(t *testing.T) {
	server, books := newTestServer(t)
	book := &models.Book{Name: "Temporary", Author: "Author", Publication: "2024"}
	if err := books.CreateBook(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	response := performRequest(server, http.MethodDelete, fmt.Sprintf("/book/%d", book.ID), "")
	if response.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", response.Code, http.StatusNoContent)
	}

	missingResponse := performRequest(server, http.MethodGet, fmt.Sprintf("/book/%d", book.ID), "")
	if missingResponse.Code != http.StatusNotFound {
		t.Fatalf("deleted book status = %d, want %d", missingResponse.Code, http.StatusNotFound)
	}
}

func performRequest(server http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
