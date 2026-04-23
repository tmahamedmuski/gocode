package main

import (
	"log"
	"net/http"

	"book-api/handler"
	"book-api/storage"
)

func main() {
	// Initialize file-based storage
	strg, err := storage.NewFileStorage("data.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize handlers
	bookHandler := handler.NewBookHandler(strg)

	// Setup routes using Go 1.22 enhanced ServeMux
	mux := http.NewServeMux()

	// REST endpoints
	mux.HandleFunc("GET /books", bookHandler.GetBooks)
	mux.HandleFunc("POST /books", bookHandler.CreateBook)
	mux.HandleFunc("GET /books/{id}", bookHandler.GetBook)
	mux.HandleFunc("PUT /books/{id}", bookHandler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", bookHandler.DeleteBook)

	// Search endpoint
	mux.HandleFunc("GET /books/search", bookHandler.SearchBooks)

	port := ":8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
