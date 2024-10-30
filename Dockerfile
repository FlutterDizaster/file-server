#syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o file-server ./cmd/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/file-server ./

CMD ["./file-server"]