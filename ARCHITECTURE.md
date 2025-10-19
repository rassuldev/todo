# System Architecture

## Overview

The Task Management System is built using a microservices architecture pattern with the following key characteristics:

- **Service Independence**: Each microservice can be developed, deployed, and scaled independently
- **Technology Agnostic**: Services can use different technologies if needed
- **Fault Isolation**: Failure in one service doesn't affect others
- **Scalability**: Each service can be scaled based on its specific load

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                          Client Layer                            │
│  (Web App, Mobile App, External Systems)                        │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ REST API / HTTP
             │
┌────────────┴────────────────────────────────────────────────────┐
│                      API Gateway (Optional)                      │
│                   Load Balancer / Ingress                        │
└──┬──────────────┬──────────────┬──────────────┬─────────────────┘
   │              │              │              │
   │ HTTP/REST    │              │              │
   │              │              │              │
┌──▼─────────┐ ┌─▼──────────┐ ┌─▼──────────┐ ┌─▼──────────────┐
│   User     │ │   Auth     │ │   Task     │ │  Notification  │
│  Service   │ │  Service   │ │  Service   │ │    Service     │
│            │ │            │ │            │ │                │
│ Port: 8081 │ │ Port: 8082 │ │ Port: 8083 │ │  Port: 8084    │
│ gRPC: 50051│ │ gRPC: 50052│ │ gRPC: 50053│ │  gRPC: 50054   │
└──┬─────────┘ └─┬──────────┘ └─┬──────────┘ └─┬──────────────┘
   │            │              │              │
   │ gRPC      │              │              │
   │           │              │              │
   └───────────┴──────────────┴──────────────┘
                      │
                      │ PostgreSQL Protocol
                      │
              ┌───────▼────────┐
              │   PostgreSQL   │
              │    Database    │
              │                │
              │ - user_db      │
              │ - task_db      │
              │ - notification_│
              │   db           │
              └────────────────┘
```

## Service Details

### 1. User Service

**Responsibility**: User management and profile operations

**Endpoints**:
- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `GET /api/users` - List users

**gRPC Methods**:
- `CreateUser`
- `GetUser`
- `UpdateUser`
- `DeleteUser`
- `ListUsers`

**Database**: `user_db`

**Tables**:
```sql
users (
  id VARCHAR(36) PRIMARY KEY,
  username VARCHAR(100) UNIQUE,
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255),
  full_name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)
```

### 2. Authorization Service

**Responsibility**: Authentication and JWT token management

**Endpoints**:
- `POST /api/auth/login` - User login
- `POST /api/auth/validate` - Validate token

**gRPC Methods**:
- `Login`
- `ValidateToken`
- `RefreshToken`
- `Logout`

**Database**: `user_db` (shared with User Service)

**Tables**:
```sql
refresh_tokens (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(36),
  token VARCHAR(500) UNIQUE,
  expires_at TIMESTAMP,
  created_at TIMESTAMP
)
```

**Security**:
- JWT tokens with HS256 signing
- Refresh token rotation
- Token expiration (24 hours for access tokens)
- Password hashing with bcrypt

### 3. Task Service

**Responsibility**: Task CRUD operations and management

**Endpoints**:
- `POST /api/tasks` - Create task
- `GET /api/tasks/{id}` - Get task
- `PUT /api/tasks/{id}` - Update task
- `DELETE /api/tasks/{id}` - Delete task
- `GET /api/tasks` - List all tasks
- `GET /api/users/{user_id}/tasks` - List user tasks

**gRPC Methods**:
- `CreateTask`
- `GetTask`
- `UpdateTask`
- `DeleteTask`
- `ListTasks`
- `ListUserTasks`

**Database**: `task_db`

**Tables**:
```sql
tasks (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255),
  description TEXT,
  status VARCHAR(50),
  priority VARCHAR(50),
  user_id VARCHAR(36),
  due_date TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)
```

**Task Status**:
- PENDING
- IN_PROGRESS
- COMPLETED
- CANCELLED

**Task Priority**:
- LOW
- MEDIUM
- HIGH
- URGENT

### 4. Notification Service

**Responsibility**: Send email and push notifications

**Endpoints**:
- `POST /api/notifications/email` - Send email
- `POST /api/notifications/push` - Send push notification

**gRPC Methods**:
- `SendEmail`
- `SendPushNotification`
- `SendTaskReminder`

**Database**: `notification_db`

**Tables**:
```sql
notifications (
  id VARCHAR(36) PRIMARY KEY,
  type VARCHAR(50),
  recipient VARCHAR(255),
  subject VARCHAR(255),
  body TEXT,
  sent BOOLEAN,
  created_at TIMESTAMP
)
```

**Notification Types**:
- EMAIL
- PUSH

## Communication Patterns

### 1. Client-to-Service Communication

**Protocol**: REST API over HTTP/HTTPS

**Why REST?**
- Universal client support
- Easy to test and debug
- Well-understood by developers
- Suitable for external clients

**Request/Response Format**: JSON

### 2. Service-to-Service Communication

**Protocol**: gRPC

**Why gRPC?**
- High performance (binary protocol)
- Strong typing with Protocol Buffers
- Bi-directional streaming support
- Built-in load balancing
- Language agnostic

### 3. Database Access

**Pattern**: Direct database access per service

**Why?**
- Simplicity
- Better performance
- Each service owns its data

## Data Management

### Database per Service Pattern

Each service has its own database to ensure:
- **Service Independence**: Services can be deployed independently
- **Technology Flexibility**: Each service can use different database technology if needed
- **Scalability**: Databases can be scaled independently
- **Fault Isolation**: Database failure affects only one service

### Data Consistency

For cross-service transactions, we use:
- **Saga Pattern**: Eventual consistency through compensating transactions
- **Event-Driven Architecture**: Services publish events for state changes

## Scalability

### Horizontal Scaling

All services are stateless and can be scaled horizontally:

```bash
# Scale to 5 replicas
kubectl scale deployment user-service --replicas=5
```

### Load Balancing

Kubernetes provides built-in load balancing:
- Service discovery
- Round-robin load balancing
- Health checks

### Database Scaling

Options for scaling the database:
- **Read Replicas**: For read-heavy workloads
- **Sharding**: For data distribution
- **Connection Pooling**: PgBouncer for connection management

## Security

### Authentication Flow

```
1. Client → POST /api/auth/login → Auth Service
2. Auth Service validates credentials with User Service (gRPC)
3. Auth Service generates JWT token
4. Client receives access_token and refresh_token
5. Client includes token in subsequent requests
6. Services validate token with Auth Service (gRPC)
```

### Security Layers

1. **Transport Security**
   - TLS/SSL for all external communication
   - mTLS for service-to-service communication (optional)

2. **Authentication**
   - JWT tokens for stateless authentication
   - Token expiration and refresh mechanism

3. **Authorization**
   - Role-based access control (RBAC)
   - Service-level permissions

4. **Data Security**
   - Password hashing with bcrypt
   - Database encryption at rest
   - Secrets management with Kubernetes Secrets

## Resilience and Reliability

### Health Checks

Each service implements:
- **Liveness Probe**: Is the service alive?
- **Readiness Probe**: Is the service ready to accept traffic?

### Circuit Breaker Pattern

Implement circuit breakers for service-to-service calls to prevent cascading failures.

### Retry Logic

Implement exponential backoff for failed requests:
```go
// Example retry logic
maxRetries := 3
backoff := time.Second

for i := 0; i < maxRetries; i++ {
    if err := makeRequest(); err == nil {
        return nil
    }
    time.Sleep(backoff * time.Duration(i+1))
}
```

### Graceful Degradation

Services should degrade gracefully when dependencies are unavailable:
- Return cached data
- Return partial results
- Provide default responses

## Monitoring and Observability

### Metrics (Prometheus)

Key metrics to monitor:
- Request rate
- Error rate
- Response time (latency)
- Resource usage (CPU, memory)

### Logging (ELK Stack)

Structured logging with:
- Log levels (DEBUG, INFO, WARN, ERROR)
- Correlation IDs for request tracing
- Contextual information

### Tracing (Jaeger)

Distributed tracing to:
- Track requests across services
- Identify bottlenecks
- Debug production issues

## Deployment Strategy

### Blue-Green Deployment

```bash
# Deploy new version (green)
kubectl apply -f k8s/deployment-v2.yaml

# Switch traffic
kubectl patch service user-service -p '{"spec":{"selector":{"version":"v2"}}}'

# Rollback if needed
kubectl patch service user-service -p '{"spec":{"selector":{"version":"v1"}}}'
```

### Canary Deployment

```bash
# Deploy canary with 10% traffic
kubectl apply -f k8s/canary-deployment.yaml

# Monitor metrics
# If successful, gradually increase traffic
# If failed, rollback
```

### Rolling Updates

Kubernetes default strategy:
- Zero downtime deployments
- Gradual rollout
- Automatic rollback on failure

## Best Practices

### 1. API Design

- RESTful principles
- Versioning (e.g., `/api/v1/users`)
- Consistent error responses
- Pagination for list endpoints
- Rate limiting

### 2. Error Handling

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

### 3. Configuration Management

- Environment variables for configuration
- ConfigMaps for non-sensitive data
- Secrets for sensitive data
- Feature flags for gradual rollouts

### 4. Documentation

- OpenAPI/Swagger for REST APIs
- Protocol Buffer documentation for gRPC
- Architecture decision records (ADRs)

## Performance Optimization

### 1. Caching

Implement caching at multiple levels:
- **Application Cache**: Redis for frequently accessed data
- **Database Cache**: PostgreSQL query cache
- **CDN**: For static assets

### 2. Database Optimization

- Proper indexing
- Query optimization
- Connection pooling
- Prepared statements

### 3. Service Optimization

- Efficient data structures
- Batch processing
- Asynchronous operations
- Compression

## Future Enhancements

### 1. API Gateway

Add Kong or custom API Gateway for:
- Rate limiting
- Request transformation
- Authentication
- Analytics

### 2. Message Queue

Implement RabbitMQ or Kafka for:
- Asynchronous processing
- Event-driven architecture
- Decoupling services

### 3. Caching Layer

Add Redis for:
- Session management
- Cache frequently accessed data
- Rate limiting counters

### 4. Service Mesh

Implement Istio for:
- Traffic management
- Security (mTLS)
- Observability
- Resilience

### 5. GraphQL Gateway

Add GraphQL layer for:
- Flexible queries
- Reduced over-fetching
- Real-time subscriptions

## Conclusion

This microservices architecture provides:
- ✅ Scalability through independent service scaling
- ✅ Reliability through fault isolation
- ✅ Maintainability through clear service boundaries
- ✅ Flexibility through independent deployment
- ✅ Performance through optimized communication protocols

The system is designed to handle production workloads with proper monitoring, security, and resilience mechanisms in place.

