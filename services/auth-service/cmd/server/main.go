package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/todo/proto/auth"
	grpcServer "github.com/todo/services/auth-service/internal/grpc"
	httpHandler "github.com/todo/services/auth-service/internal/http"
	"github.com/todo/services/auth-service/internal/jwt"
	"github.com/todo/services/auth-service/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	// Get environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "user_db")
	grpcPort := getEnv("GRPC_PORT", "50052")
	httpPort := getEnv("HTTP_PORT", "8082")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	repo, err := repository.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Create JWT manager
	jwtManager := jwt.NewJWTManager(jwtSecret, 24*time.Hour)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterAuthServiceServer(s, grpcServer.NewAuthServer(repo, jwtManager))

		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	router := mux.NewRouter()
	handler := httpHandler.NewHandler(repo, jwtManager)
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
