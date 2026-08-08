FROM golang:1.26.3-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]