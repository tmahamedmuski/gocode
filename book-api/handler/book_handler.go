package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"book-api/model"
	"book-api/storage"
)

type BookHandler struct {
	Store storage.Storage
}

func NewBookHandler(store storage.Storage) *BookHandler {
	return &BookHandler{Store: store}
}

// WriteJSON is a helper for JSON responses
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// WriteError is a helper for error responses
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// CreateBook new ones (POST /books)
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if book.BookID == "" {
		WriteError(w, http.StatusBadRequest, "bookId is required")
		return
	}

	if err := h.Store.Create(book); err != nil {
		if err == storage.ErrBookExists {
			WriteError(w, http.StatusConflict, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	WriteJSON(w, http.StatusCreated, book)
}

// GetBooks GET /books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Default limit
	offset := 0

	query := r.URL.Query()
	if l := query.Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if o := query.Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	books, total, err := h.Store.GetAll(offset, limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	response := map[string]interface{}{
		"data":   books,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	WriteJSON(w, http.StatusOK, response)
}

// GetBook GET /books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "ID is required")
		return
	}

	book, err := h.Store.GetByID(id)
	if err != nil {
		if err == storage.ErrBookNotFound {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve book")
		return
	}

	WriteJSON(w, http.StatusOK, book)
}

// UpdateBook PUT /books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "ID is required")
		return
	}

	// Fetch existing book first
	book, err := h.Store.GetByID(id)
	if err != nil {
		if err == storage.ErrBookNotFound {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve book for update")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.Store.Update(id, book); err != nil {
		if err == storage.ErrBookNotFound {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update book")
		return
	}

	// Fetch the updated book to return
	updatedBook, _ := h.Store.GetByID(id)
	WriteJSON(w, http.StatusOK, updatedBook)
}

// DeleteBook DELETE /books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "ID is required")
		return
	}

	if err := h.Store.Delete(id); err != nil {
		if err == storage.ErrBookNotFound {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchBooks GET /books/search?q=<keyword>
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		WriteError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	allBooks := h.Store.GetAllBooks()
	if len(allBooks) == 0 {
		WriteJSON(w, http.StatusOK, []model.Book{})
		return
	}

	keyword := strings.ToLower(query)

	// Optimization Approach: Split into Goroutines
	numWorkers := 4
	booksPerWorker := int(math.Ceil(float64(len(allBooks)) / float64(numWorkers)))

	resultsChan := make(chan []model.Book, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		start := i * booksPerWorker
		end := start + booksPerWorker
		if start >= len(allBooks) {
			break
		}
		if end > len(allBooks) {
			end = len(allBooks)
		}

		wg.Add(1)
		go func(booksSlice []model.Book) {
			defer wg.Done()
			var localMatches []model.Book
			for _, b := range booksSlice {
				if strings.Contains(strings.ToLower(b.Title), keyword) ||
					strings.Contains(strings.ToLower(b.Description), keyword) {
					localMatches = append(localMatches, b)
				}
			}
			resultsChan <- localMatches
		}(allBooks[start:end])
	}

	// Close the channel once all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var finalMatches []model.Book
	for partialMatches := range resultsChan {
		finalMatches = append(finalMatches, partialMatches...)
	}

	if finalMatches == nil {
		finalMatches = []model.Book{}
	}

	WriteJSON(w, http.StatusOK, finalMatches)
}
