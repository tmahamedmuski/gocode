package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"book-api/model"
)

var (
	ErrBookNotFound = errors.New("book not found")
	ErrBookExists   = errors.New("book already exists")
)

type Storage interface {
	GetAll(offset, limit int) ([]model.Book, int, error)
	GetByID(id string) (model.Book, error)
	Create(book model.Book) error
	Update(id string, book model.Book) error
	Delete(id string) error
	GetAllBooks() []model.Book
}

type FileStorage struct {
	mu       sync.RWMutex
	filename string
	books    map[string]model.Book
	ordered  []string // keeps track of ordering for consistent pagination
}

func NewFileStorage(filename string) (*FileStorage, error) {
	fs := &FileStorage{
		filename: filename,
		books:    make(map[string]model.Book),
		ordered:  make([]string, 0),
	}
	err := fs.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStorage) load() error {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	var bookList []model.Book
	if err := json.Unmarshal(data, &bookList); err != nil {
		return err
	}
	for _, book := range bookList {
		fs.books[book.BookID] = book
		fs.ordered = append(fs.ordered, book.BookID)
	}
	return nil
}

func (fs *FileStorage) save() error {
	var bookList []model.Book
	for _, id := range fs.ordered {
		bookList = append(bookList, fs.books[id])
	}
	data, err := json.MarshalIndent(bookList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) GetAll(offset, limit int) ([]model.Book, int, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	total := len(fs.ordered)
	if offset >= total {
		return []model.Book{}, total, nil
	}

	end := offset + limit
	if end > total || limit == 0 {
		end = total // return all if limit is 0
	}

	var result []model.Book
	for _, id := range fs.ordered[offset:end] {
		result = append(result, fs.books[id])
	}
	return result, total, nil
}

func (fs *FileStorage) GetByID(id string) (model.Book, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	book, exists := fs.books[id]
	if !exists {
		return model.Book{}, ErrBookNotFound
	}
	return book, nil
}

func (fs *FileStorage) Create(book model.Book) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.books[book.BookID]; exists {
		return ErrBookExists
	}

	fs.books[book.BookID] = book
	fs.ordered = append(fs.ordered, book.BookID)

	return fs.save()
}

func (fs *FileStorage) Update(id string, book model.Book) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.books[id]; !exists {
		return ErrBookNotFound
	}
	
	// Enforce ID matches
	book.BookID = id
	fs.books[id] = book

	return fs.save()
}

func (fs *FileStorage) Delete(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.books[id]; !exists {
		return ErrBookNotFound
	}

	delete(fs.books, id)

	// Remove from ordered slice
	for i, existingID := range fs.ordered {
		if existingID == id {
			fs.ordered = append(fs.ordered[:i], fs.ordered[i+1:]...)
			break
		}
	}

	return fs.save()
}

// GetAllBooks returns a copy of all books for searching
func (fs *FileStorage) GetAllBooks() []model.Book {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var result []model.Book
	for _, id := range fs.ordered {
		result = append(result, fs.books[id])
	}
	return result
}
