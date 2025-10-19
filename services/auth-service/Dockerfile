# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

# Download dependencies
RUN cd services/auth-service && go mod download

# Copy source code
COPY . .

# Build the application
RUN cd services/auth-service && CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth-service .

EXPOSE 50052 8082

CMD ["./auth-service"]

