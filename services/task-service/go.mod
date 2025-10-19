module github.com/todo/services/task-service

go 1.21

replace github.com/todo => ../..

require (
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/todo v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

