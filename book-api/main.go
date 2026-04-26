package main

import (
	"log"
	"net/http"

	"book-api/handler"
	"book-api/storage"
)

func main() {

	strg, err := storage.NewFileStorage("data.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	bookHandler := handler.NewBookHandler(strg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", bookHandler.GetBooks)
	mux.HandleFunc("POST /books", bookHandler.CreateBook)
	mux.HandleFunc("GET /books/{id}", bookHandler.GetBook)
	mux.HandleFunc("PUT /books/{id}", bookHandler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", bookHandler.DeleteBook)

	mux.HandleFunc("GET /books/search", bookHandler.SearchBooks)

	port := ":8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
