# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/task-service/go.mod services/task-service/go.sum ./services/task-service/

# Download dependencies
RUN cd services/task-service && go mod download

# Copy source code
COPY . .

# Build the application
RUN cd services/task-service && CGO_ENABLED=0 GOOS=linux go build -o /app/task-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/task-service .

EXPOSE 50053 8083

CMD ["./task-service"]

