package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"book-api/model"
	"book-api/storage"
)

func setupTestStorage() storage.Storage {
	strg, _ := storage.NewFileStorage("test_data.json")
	return strg
}

func cleanupTestStorage() {
	os.Remove("test_data.json")
}

func TestCreateBook(t *testing.T) {
	strg := setupTestStorage()
	defer cleanupTestStorage()

	handler := NewBookHandler(strg)

	book := model.Book{
		BookID: "1",
		Title:  "Test Book",
	}

	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBook(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %v", res.StatusCode)
	}

	var createdBook model.Book
	json.NewDecoder(res.Body).Decode(&createdBook)
	if createdBook.BookID != "1" {
		t.Errorf("Expected book ID '1', got '%s'", createdBook.BookID)
	}
}

func TestGetBooks(t *testing.T) {
	strg := setupTestStorage()
	defer cleanupTestStorage()

	strg.Create(model.Book{BookID: "1", Title: "Book 1"})
	strg.Create(model.Book{BookID: "2", Title: "Book 2"})

	handler := NewBookHandler(strg)

	req := httptest.NewRequest(http.MethodGet, "/books?limit=1", nil)
	w := httptest.NewRecorder()

	handler.GetBooks(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", res.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)

	data := response["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("Expected 1 book due to limit, got %d", len(data))
	}
}
