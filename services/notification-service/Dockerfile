# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/

# Download dependencies
RUN cd services/notification-service && go mod download

# Copy source code
COPY . .

# Build the application
RUN cd services/notification-service && CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/notification-service .

EXPOSE 50054 8084

CMD ["./notification-service"]

