# Stage 1: Build the Go app
FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ingestion-service ./cmd/main.go

# Stage 2: Create lightweight image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/ingestion-service .

# Load .env file (optional)
# COPY .env .env

CMD ["./ingestion-service"]
