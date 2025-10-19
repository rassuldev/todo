# Task Management System (To-Do List) - Microservices Architecture

A robust microservices-based task management system built with Go, featuring user management, authentication, task operations, and notifications. The system uses gRPC for inter-service communication and REST API for client-facing endpoints.

## 🏗️ Architecture

The system consists of four main microservices:

1. **User Service** - User registration, profiles, and management
2. **Authorization Service** - Authentication with JWT tokens and access control
3. **Task Service** - CRUD operations for tasks
4. **Notification Service** - Email and push notifications

## 🛠️ Technologies

- **Language**: Go 1.21
- **Communication**: gRPC (internal), REST API (external)
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **API Router**: Gorilla Mux

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Kubernetes (kubectl) and Minikube/Kind (for local K8s deployment)
- PostgreSQL 15
- Protocol Buffers compiler (protoc)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd todo
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Generate Protocol Buffers

Install the required tools:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Generate proto files:

```bash
make proto
```

### 4. Run with Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 5. Run Locally (Development)

Start PostgreSQL:

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_db \
  -p 5432:5432 \
  postgres:15-alpine
```

Run each service in separate terminals:

```bash
# Terminal 1 - User Service
cd services/user-service
go run cmd/server/main.go

# Terminal 2 - Auth Service
cd services/auth-service
go run cmd/server/main.go

# Terminal 3 - Task Service
cd services/task-service
go run cmd/server/main.go

# Terminal 4 - Notification Service
cd services/notification-service
go run cmd/server/main.go
```

## ☸️ Kubernetes Deployment

### Build Docker Images

```bash
make docker-build
```

### Deploy to Kubernetes

```bash
# Start Minikube (if using Minikube)
minikube start

# Deploy all services
make k8s-deploy

# Check deployments
kubectl get pods
kubectl get services

# Access services (if using Minikube)
minikube service user-service-http
minikube service auth-service-http
minikube service task-service-http
minikube service notification-service-http

# Delete deployments
make k8s-delete
```

## 📡 API Endpoints

### User Service (Port 8081)

#### Create User
```bash
POST /api/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password",
  "full_name": "John Doe"
}
```

#### Get User
```bash
GET /api/users/{id}
```

#### Update User
```bash
PUT /api/users/{id}
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "full_name": "John Doe Updated"
}
```

#### Delete User
```bash
DELETE /api/users/{id}
```

#### List Users
```bash
GET /api/users?page=1&page_size=10
```

### Auth Service (Port 8082)

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password"
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "uuid-refresh-token",
  "expires_at": "2024-01-01T00:00:00Z"
}
```

#### Validate Token
```bash
POST /api/auth/validate
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "valid": true,
  "user_id": "user-uuid",
  "username": "john_doe"
}
```

### Task Service (Port 8083)

#### Create Task
```bash
POST /api/tasks
Content-Type: application/json

{
  "title": "Complete project",
  "description": "Finish the microservices implementation",
  "priority": "HIGH",
  "user_id": "user-uuid",
  "due_date": "2024-12-31T23:59:59Z"
}
```

#### Get Task
```bash
GET /api/tasks/{id}
```

#### Update Task
```bash
PUT /api/tasks/{id}
Content-Type: application/json

{
  "title": "Complete project",
  "description": "Finish the microservices implementation",
  "status": "IN_PROGRESS",
  "priority": "HIGH",
  "due_date": "2024-12-31T23:59:59Z"
}
```

#### Delete Task
```bash
DELETE /api/tasks/{id}
```

#### List All Tasks
```bash
GET /api/tasks?page=1&page_size=10
```

#### List User Tasks
```bash
GET /api/users/{user_id}/tasks?page=1&page_size=10&status=PENDING
```

### Notification Service (Port 8084)

#### Send Email
```bash
POST /api/notifications/email
Content-Type: application/json

{
  "to": "user@example.com",
  "subject": "Task Reminder",
  "body": "Don't forget about your task!"
}
```

#### Send Push Notification
```bash
POST /api/notifications/push
Content-Type: application/json

{
  "device_token": "device-token",
  "title": "Task Reminder",
  "body": "Don't forget about your task!"
}
```

## 🧪 Testing Examples

### Complete User Flow

```bash
# 1. Create a user
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# 2. Login
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# Save the access_token from response

# 3. Create a task
curl -X POST http://localhost:8083/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Task",
    "description": "This is a test task",
    "priority": "MEDIUM",
    "user_id": "user-uuid-from-step-1"
  }'

# 4. List user tasks
curl -X GET "http://localhost:8083/api/users/{user-id}/tasks?page=1&page_size=10"

# 5. Send notification
curl -X POST http://localhost:8084/api/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Task Created",
    "body": "Your task has been created successfully!"
  }'
```

## 🗂️ Project Structure

```
todo/
├── proto/                          # Protocol Buffer definitions
│   ├── user/
│   ├── auth/
│   ├── task/
│   └── notification/
├── services/
│   ├── user-service/              # User management service
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── models/
│   │   │   ├── repository/
│   │   │   ├── grpc/
│   │   │   └── http/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── auth-service/              # Authentication service
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── models/
│   │   │   ├── jwt/
│   │   │   ├── repository/
│   │   │   ├── grpc/
│   │   │   └── http/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── task-service/              # Task management service
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── models/
│   │   │   ├── repository/
│   │   │   ├── grpc/
│   │   │   └── http/
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── notification-service/      # Notification service
│       ├── cmd/server/
│       ├── internal/
│       │   ├── models/
│       │   ├── email/
│       │   ├── push/
│       │   ├── repository/
│       │   ├── grpc/
│       │   └── http/
│       ├── Dockerfile
│       └── go.mod
├── k8s/                           # Kubernetes configurations
│   ├── postgres-deployment.yaml
│   ├── user-service-deployment.yaml
│   ├── auth-service-deployment.yaml
│   ├── task-service-deployment.yaml
│   └── notification-service-deployment.yaml
├── docker-compose.yaml            # Docker Compose configuration
├── Makefile                       # Build automation
├── go.mod                         # Go module file
└── README.md                      # This file
```

## 🔧 Configuration

### Environment Variables

Each service can be configured using environment variables:

#### User Service
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: user_db)
- `GRPC_PORT` - gRPC port (default: 50051)
- `HTTP_PORT` - HTTP port (default: 8081)

#### Auth Service
- Same as User Service, plus:
- `JWT_SECRET` - Secret key for JWT signing
- `GRPC_PORT` - gRPC port (default: 50052)
- `HTTP_PORT` - HTTP port (default: 8082)

#### Task Service
- Same database configs
- `DB_NAME` - Database name (default: task_db)
- `GRPC_PORT` - gRPC port (default: 50053)
- `HTTP_PORT` - HTTP port (default: 8083)

#### Notification Service
- Same database configs
- `DB_NAME` - Database name (default: notification_db)
- `GRPC_PORT` - gRPC port (default: 50054)
- `HTTP_PORT` - HTTP port (default: 8084)
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `SMTP_FROM` - From email address
- `PUSH_API_KEY` - Push notification API key

## 🔐 Security Considerations

1. **JWT Secret**: Change the default JWT secret in production
2. **Database Passwords**: Use strong passwords and secrets management
3. **HTTPS**: Use TLS/SSL in production
4. **API Gateway**: Consider adding an API Gateway for production
5. **Rate Limiting**: Implement rate limiting for public endpoints
6. **Input Validation**: All inputs are validated at service level

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs todo-postgres
```

### Service Not Starting
```bash
# Check service logs
docker-compose logs user-service
docker-compose logs auth-service
docker-compose logs task-service
docker-compose logs notification-service
```

### Port Already in Use
```bash
# Find process using port
lsof -i :8081

# Kill process
kill -9 <PID>
```

## 📚 Additional Resources

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Go Documentation](https://golang.org/doc/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

Built as a demonstration of microservices architecture in Go.

## 🙏 Acknowledgments

- Go community for excellent libraries
- gRPC team for the communication framework
- Docker and Kubernetes communities

