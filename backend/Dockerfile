# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o ./backend main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/backend .
COPY --from=builder /app/docs ./docs
COPY backend-config.yaml /configs/backend-config.yaml

CMD ["./backend"]