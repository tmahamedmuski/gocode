# Book API (Go)

A high-performance REST API built with Go for managing a collection of books.
This project supports CRUD operations, pagination, concurrent search, file-based storage, and containerized deployment.


### Core Features

* Create, Read, Update, Delete (CRUD) books
* RESTful API using Go’s native `net/http`
* JSON request/response handling

### Performance Optimizations

* **O(1) lookups** using in-memory map
* Ordered slice for efficient pagination
* Thread-safe operations with `sync.RWMutex`

### Concurrent Search

* Search endpoint uses **goroutines**
* Dataset split into chunks for parallel processing
* Case-insensitive matching on title and description
* Results combined using channels and `sync.WaitGroup`

### Storage(not recomanded for production)

* File-based persistence using `data.json`
* Auto-created file if not present
* Data written on every mutation

### Pagination

* Supports `offset` and `limit` parameters

### Testing

* Unit tests using `net/http/httptest`

### Deployment Ready

* Docker support
* Kubernetes manifests included



## Tech Stack

* **Language:** Go 1.22
* **Routing:** `http.ServeMux` (native Go router)
* **Storage:** JSON file
* **Concurrency:** Goroutines, Channels
* **Containerization:** Docker
* **Orchestration:** Kubernetes



## Project Structure
book-api/
├── main.go             
├── model/
│   └── book.go         
├── storage/
│   └── file_storage.go 
├── handler/
│   └── book_handler.go 
├── k8s/
│   ├── deployment.yaml
│   └── service.yaml
├── Dockerfile
├── data.json
└── README.md

## API Endpoints

| Method | Endpoint           | Description                      |
| ------ | ------------------ | -------------------------------- |
| POST   | `/books`           | Create a new book                |
| GET    | `/books`           | List all books (with pagination) |
| GET    | `/books/{id}`      | Get book by ID                   |
| PUT    | `/books/{id}`      | Update book                      |
| DELETE | `/books/{id}`      | Delete book                      |
| GET    | `/books/search?q=` | Search books                     |



##  Running the Application

Start the server:
go run main.go

Server runs on:
http://localhost:8080


### Create a Book

bash
curl -X POST http://localhost:8080/books \
-H "Content-Type: application/json" \
-d '{"bookId":"1","title":"Go Programming","authorId":"A1","pages":250,"price":29.99}'


### Get Books (Pagination)


curl "http://localhost:8080/books?offset=0&limit=5"


### Search Books


curl "http://localhost:8080/books/search?q=go"


### Update Book


curl -X PUT http://localhost:8080/books/1 \
-H "Content-Type: application/json" \
-d '{"title":"Advanced Go","pages":300}'


### Delete Book


curl -X DELETE http://localhost:8080/books/1


##  Running Tests


go test ./... -v




## Docker

Build the image:

docker build -t book-api .


Run the container:


docker run -p 8080:8080 book-api


##  Kubernetes

Apply deployment and service:


kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml


Access the service using your cluster tools 



## Highlights

* Clean architecture (Handler → Storage → Model)
* Native Go routing (no external frameworks)
* Efficient concurrent search implementation
* Production-ready with Docker & Kubernetes



## Conclusion

This project was developed through my own hands-on coding and learning process using Go. While building the system, I encountered various errors and challenges, where I used ChatGPT as a guidance tool to debug issues and clarify concepts. The complete implementation was not directly generated, but rather designed, coded, and refined by me. This experience significantly improved my understanding of REST API development, concurrency, and practical backend engineering.