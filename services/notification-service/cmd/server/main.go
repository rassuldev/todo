package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	pb "github.com/todo/proto/notification"
	"github.com/todo/services/notification-service/internal/email"
	grpcServer "github.com/todo/services/notification-service/internal/grpc"
	httpHandler "github.com/todo/services/notification-service/internal/http"
	"github.com/todo/services/notification-service/internal/push"
	"github.com/todo/services/notification-service/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	// Get environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "notification_db")
	grpcPort := getEnv("GRPC_PORT", "50054")
	httpPort := getEnv("HTTP_PORT", "8084")
	pushAPIKey := getEnv("PUSH_API_KEY", "")

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	repo, err := repository.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Create email and push senders
	emailSender := email.NewEmailSender()
	pushSender := push.NewPushSender(pushAPIKey)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterNotificationServiceServer(s, grpcServer.NewNotificationServer(repo, emailSender, pushSender))

		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	router := mux.NewRouter()
	handler := httpHandler.NewHandler(repo, emailSender, pushSender)
	handler.RegisterRoutes(router)

	log.Printf("HTTP server listening on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, router); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
