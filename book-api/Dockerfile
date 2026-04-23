FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o book-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/book-api .

EXPOSE 8080

CMD ["./book-api"]
