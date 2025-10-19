# Deployment Guide

This guide provides detailed instructions for deploying the Task Management System in different environments.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Compose](#docker-compose)
3. [Kubernetes](#kubernetes)
4. [Production Considerations](#production-considerations)

## Local Development

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Protocol Buffers compiler

### Setup Steps

1. **Install Go Dependencies**
   ```bash
   go mod download
   ```

2. **Generate Protocol Buffers**
   ```bash
   # Install protoc plugins
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   
   # Generate proto files
   make proto
   ```

3. **Start PostgreSQL**
   ```bash
   docker run -d \
     --name postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=user_db \
     -p 5432:5432 \
     postgres:15-alpine
   ```

4. **Run Services**
   
   Open 4 terminal windows and run:
   
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

## Docker Compose

### Quick Start

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f user-service

# Stop all services
docker-compose down

# Remove volumes (clean start)
docker-compose down -v
```

### Rebuild After Changes

```bash
# Rebuild specific service
docker-compose build user-service

# Rebuild all services
docker-compose build

# Rebuild and restart
docker-compose up -d --build
```

### Accessing Services

- User Service HTTP: http://localhost:8081
- Auth Service HTTP: http://localhost:8082
- Task Service HTTP: http://localhost:8083
- Notification Service HTTP: http://localhost:8084
- PostgreSQL: localhost:5432

## Kubernetes

### Local Kubernetes (Minikube)

1. **Start Minikube**
   ```bash
   minikube start --cpus=4 --memory=8192
   
   # Enable Minikube to use local Docker images
   eval $(minikube docker-env)
   ```

2. **Build Docker Images**
   ```bash
   make docker-build
   ```

3. **Deploy to Kubernetes**
   ```bash
   # Deploy all resources
   kubectl apply -f k8s/
   
   # Or use Makefile
   make k8s-deploy
   ```

4. **Check Deployment Status**
   ```bash
   # Check pods
   kubectl get pods
   
   # Check services
   kubectl get services
   
   # Check deployments
   kubectl get deployments
   
   # Check logs
   kubectl logs -f deployment/user-service
   ```

5. **Access Services**
   ```bash
   # Get service URLs
   minikube service user-service-http --url
   minikube service auth-service-http --url
   minikube service task-service-http --url
   minikube service notification-service-http --url
   
   # Or open in browser
   minikube service user-service-http
   ```

6. **Scale Services**
   ```bash
   # Scale up
   kubectl scale deployment user-service --replicas=3
   
   # Scale down
   kubectl scale deployment user-service --replicas=1
   ```

7. **Update Deployment**
   ```bash
   # After making changes and rebuilding image
   kubectl rollout restart deployment/user-service
   
   # Check rollout status
   kubectl rollout status deployment/user-service
   ```

8. **Delete Deployment**
   ```bash
   # Delete all resources
   kubectl delete -f k8s/
   
   # Or use Makefile
   make k8s-delete
   ```

### Production Kubernetes (EKS, GKE, AKS)

1. **Push Docker Images to Registry**
   ```bash
   # Tag images
   docker tag todo-user-service:latest your-registry/todo-user-service:v1.0.0
   docker tag todo-auth-service:latest your-registry/todo-auth-service:v1.0.0
   docker tag todo-task-service:latest your-registry/todo-task-service:v1.0.0
   docker tag todo-notification-service:latest your-registry/todo-notification-service:v1.0.0
   
   # Push images
   docker push your-registry/todo-user-service:v1.0.0
   docker push your-registry/todo-auth-service:v1.0.0
   docker push your-registry/todo-task-service:v1.0.0
   docker push your-registry/todo-notification-service:v1.0.0
   ```

2. **Update Kubernetes Manifests**
   
   Update image references in deployment files:
   ```yaml
   image: your-registry/todo-user-service:v1.0.0
   imagePullPolicy: Always
   ```

3. **Create Secrets**
   ```bash
   # Database credentials
   kubectl create secret generic db-credentials \
     --from-literal=username=postgres \
     --from-literal=password=your-secure-password
   
   # JWT secret
   kubectl create secret generic jwt-secret \
     --from-literal=secret=your-jwt-secret-key
   
   # SMTP credentials
   kubectl create secret generic smtp-credentials \
     --from-literal=username=your-smtp-username \
     --from-literal=password=your-smtp-password
   ```

4. **Deploy**
   ```bash
   kubectl apply -f k8s/
   ```

5. **Set Up Ingress**
   ```bash
   # Install ingress controller (nginx example)
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
   
   # Create ingress resource
   kubectl apply -f k8s/ingress.yaml
   ```

## Production Considerations

### 1. Security

- **Use Secrets Management**: Store sensitive data in Kubernetes Secrets or external secret managers (AWS Secrets Manager, HashiCorp Vault)
- **Enable TLS**: Use cert-manager for automatic certificate management
- **Network Policies**: Restrict inter-pod communication
- **RBAC**: Implement proper Role-Based Access Control
- **Security Scanning**: Scan Docker images for vulnerabilities

### 2. Database

- **Managed Database**: Use managed PostgreSQL (AWS RDS, GCP Cloud SQL, Azure Database)
- **Backups**: Configure automatic backups
- **High Availability**: Set up replication and failover
- **Connection Pooling**: Implement PgBouncer for connection pooling
- **Monitoring**: Set up database monitoring and alerts

### 3. Monitoring and Logging

- **Prometheus + Grafana**: Metrics collection and visualization
- **ELK Stack**: Centralized logging (Elasticsearch, Logstash, Kibana)
- **Jaeger**: Distributed tracing
- **Health Checks**: Implement liveness and readiness probes

Example health check:
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8081
  initialDelaySeconds: 30
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /ready
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 5
```

### 4. Performance

- **Horizontal Pod Autoscaling**: Automatically scale based on metrics
  ```bash
  kubectl autoscale deployment user-service --cpu-percent=70 --min=2 --max=10
  ```
- **Resource Limits**: Set CPU and memory limits
  ```yaml
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "500m"
  ```
- **Caching**: Implement Redis for caching
- **CDN**: Use CDN for static assets

### 5. CI/CD

Example GitHub Actions workflow:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Build Docker images
        run: make docker-build
      
      - name: Push to registry
        run: |
          docker push your-registry/todo-user-service:latest
          # ... other services
      
      - name: Deploy to Kubernetes
        run: kubectl apply -f k8s/
```

### 6. API Gateway

Consider adding an API Gateway (Kong, Ambassador, or custom) for:
- Rate limiting
- Authentication
- Request routing
- API versioning
- Analytics

### 7. Service Mesh

For complex microservices communication, consider:
- Istio
- Linkerd
- Consul

Benefits:
- Traffic management
- Security (mTLS)
- Observability
- Resilience

## Troubleshooting

### Common Issues

1. **Pods not starting**
   ```bash
   kubectl describe pod <pod-name>
   kubectl logs <pod-name>
   ```

2. **Database connection failed**
   ```bash
   # Check if PostgreSQL is running
   kubectl get pods | grep postgres
   
   # Check service
   kubectl get svc postgres-service
   
   # Test connection from pod
   kubectl exec -it <pod-name> -- nc -zv postgres-service 5432
   ```

3. **Image pull errors**
   ```bash
   # Check image name and tag
   kubectl describe pod <pod-name>
   
   # For local development with Minikube
   eval $(minikube docker-env)
   make docker-build
   ```

4. **Out of resources**
   ```bash
   # Check node resources
   kubectl top nodes
   
   # Check pod resources
   kubectl top pods
   ```

### Health Checks

```bash
# Check if services are responding
curl http://localhost:8081/api/users
curl http://localhost:8082/api/auth/login
curl http://localhost:8083/api/tasks
curl http://localhost:8084/api/notifications/email
```

## Maintenance

### Database Migrations

For schema changes, consider using migration tools:
- golang-migrate
- goose
- sql-migrate

### Rolling Updates

```bash
# Update image
kubectl set image deployment/user-service user-service=your-registry/todo-user-service:v1.1.0

# Monitor rollout
kubectl rollout status deployment/user-service

# Rollback if needed
kubectl rollout undo deployment/user-service
```

### Backup and Restore

```bash
# Backup PostgreSQL
kubectl exec postgres-0 -- pg_dump -U postgres user_db > backup.sql

# Restore
kubectl exec -i postgres-0 -- psql -U postgres user_db < backup.sql
```

## Support

For issues and questions:
1. Check logs: `kubectl logs -f <pod-name>`
2. Describe resources: `kubectl describe <resource> <name>`
3. Review events: `kubectl get events --sort-by='.lastTimestamp'`

