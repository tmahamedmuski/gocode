# Book API - Go Assignment

This project is a REST API written in Go for managing a collection of Books. It addresses all the core requirements as well as the bonus points (Pagination, Docker, Kubernetes, Unit Testing).

## Features & Optimizations

### 1. Robust File-Based Storage & Indexing Mechanism
The storage layer (`storage/file_storage.go`) uses `data.json` as the persistent backend, but **optimizes reads and writes** by implementing a thread-safe in-memory mapping:
- Lookups (`GET /books/{id}`) are **O(1)** due to map lookup.
- Fast, ordered lookups (`GET /books?limit=X&offset=Y`) are maintained using an ordered slice of IDs, efficiently supporting pagination.
- Durability is guaranteed by a `sync.RWMutex`. The JSON file is re-written on mutation (`POST`, `PUT`, `DELETE`).

### 2. Search Optimization (Concurrency)
The `/books/search?q=<keyword>` endpoint executes highly-optimized queries using Goroutines and Channels:
- The entire dataset is split evenly among several **goroutines** (workers).
- Case-insensitive substring matching is performed in parallel for `title` and `description`.
- A `sync.WaitGroup` ensures all chunks are matched, and a channel gathers local partial matches into a unified matching response list.

### 3. Native Routing
Implemented using Go 1.22's newly enhanced standard library `http.ServeMux` enabling parameterized endpoints like `GET /books/{id}` natively, without third party routers.

## Running Locally

Requirements: Go 1.22+

1. Start the server (will automatically create `data.json`):
   ```bash
   go run main.go
   ```
2. The server will start on port `8080`.

**Example Endpoints:**
- `GET http://localhost:8080/books`
- `GET http://localhost:8080/books?offset=0&limit=5`
- `POST http://localhost:8080/books`
- `GET http://localhost:8080/books/search?q=novel`

## Testing

A suite of unit tests has been implemented to test Handler integrations and concurrency mapping using `net/http/httptest`.

```bash
go test ./... -v
```

## Docker Containerization

To run this application as a container:

```bash
docker build -t book-api:latest .
docker run -p 8080:8080 book-api:latest
```

## Kubernetes Deployment (Local with Minikube / Kind)

Manifests are located in the `k8s/` directory. By default, it looks for the image `book-api:latest` locally with `imagePullPolicy: Never`.

1. Ensure Minikube is started:
   ```bash
   minikube start
   ```

2. Point your Docker CLI to Minikube's Docker daemon, then build the image so K8s can see it:
   ```bash
   eval $(minikube docker-env)
   docker build -t book-api:latest .
   ```

3. Apply the manifests:
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   ```

4. Expose the service locally (Minikube):
   ```bash
   minikube service book-api-service
   ```
