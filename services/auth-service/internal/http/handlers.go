package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todo/services/auth-service/internal/jwt"
	"github.com/todo/services/auth-service/internal/repository"
)

type Handler struct {
	repo       *repository.PostgresRepository
	jwtManager *jwt.JWTManager
}

func NewHandler(repo *repository.PostgresRepository, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid    bool   `json:"valid"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate credentials
	userID, username, err := h.repo.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, expiresAt, err := h.jwtManager.Generate(userID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := h.jwtManager.Validate(req.Token)
	if err != nil {
		response := ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ValidateResponse{
		Valid:    true,
		UserID:   claims.UserID,
		Username: claims.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/api/auth/validate", h.Validate).Methods("POST")
}
