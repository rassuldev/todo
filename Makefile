.PHONY: proto clean build docker-build k8s-deploy k8s-delete

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/task/task.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/notification/notification.proto

# Build all services
build:
	@echo "Building services..."
	@cd services/user-service && go build -o ../../bin/user-service ./cmd/server
	@cd services/auth-service && go build -o ../../bin/auth-service ./cmd/server
	@cd services/task-service && go build -o ../../bin/task-service ./cmd/server
	@cd services/notification-service && go build -o ../../bin/notification-service ./cmd/server

# Build Docker images
docker-build:
	@echo "Building Docker images..."
	@docker build -t todo-user-service:latest -f services/user-service/Dockerfile .
	@docker build -t todo-auth-service:latest -f services/auth-service/Dockerfile .
	@docker build -t todo-task-service:latest -f services/task-service/Dockerfile .
	@docker build -t todo-notification-service:latest -f services/notification-service/Dockerfile .

# Start all services with docker-compose
up:
	@docker-compose up -d

# Stop all services
down:
	@docker-compose down

# Deploy to Kubernetes
k8s-deploy:
	@kubectl apply -f k8s/

# Delete from Kubernetes
k8s-delete:
	@kubectl delete -f k8s/

# Clean build artifacts
clean:
	@rm -rf bin/
	@rm -f proto/**/*.pb.go

# Run tests
test:
	@go test ./...

