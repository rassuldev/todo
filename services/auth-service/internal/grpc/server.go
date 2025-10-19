package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/todo/proto/auth"
	"github.com/todo/services/auth-service/internal/jwt"
	"github.com/todo/services/auth-service/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	repo       *repository.PostgresRepository
	jwtManager *jwt.JWTManager
}

func NewAuthServer(repo *repository.PostgresRepository, jwtManager *jwt.JWTManager) *AuthServer {
	return &AuthServer{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate credentials
	userID, username, err := s.repo.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Error: err.Error(),
		}, nil
	}

	// Generate access token
	accessToken, expiresAt, err := s.jwtManager.Generate(userID, username)
	if err != nil {
		return &pb.LoginResponse{
			Error: err.Error(),
		}, nil
	}

	// Generate refresh token
	refreshToken := uuid.New().String()
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	// Store refresh token
	if err := s.repo.StoreRefreshToken(userID, refreshToken, refreshExpiresAt); err != nil {
		return &pb.LoginResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := s.jwtManager.Validate(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:    true,
		UserId:   claims.UserID,
		Username: claims.Username,
	}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Validate refresh token
	userID, err := s.repo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Error: err.Error(),
		}, nil
	}

	// Get username from user service (simplified - in real app would call user service via gRPC)
	username := userID // Placeholder

	// Generate new access token
	accessToken, expiresAt, err := s.jwtManager.Generate(userID, username)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Error: err.Error(),
		}, nil
	}

	// Generate new refresh token
	newRefreshToken := uuid.New().String()
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Delete old refresh token
	s.repo.DeleteRefreshToken(req.RefreshToken)

	// Store new refresh token
	if err := s.repo.StoreRefreshToken(userID, newRefreshToken, refreshExpiresAt); err != nil {
		return &pb.RefreshTokenResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// Validate and get claims from token
	_, err := s.jwtManager.Validate(req.Token)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// In a real implementation, you would add the token to a blacklist
	// For now, we just return success
	return &pb.LogoutResponse{
		Success: true,
	}, nil
}
