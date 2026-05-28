# Stage 1: build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o transit-ni .

# Stage 2: run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/transit-ni .

EXPOSE 8081

CMD ["./transit-ni"]
