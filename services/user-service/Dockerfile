# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/

# Download dependencies
RUN cd services/user-service && go mod download

# Copy source code
COPY . .

# Build the application
RUN cd services/user-service && CGO_ENABLED=0 GOOS=linux go build -o /app/user-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/user-service .

EXPOSE 50051 8081

CMD ["./user-service"]

