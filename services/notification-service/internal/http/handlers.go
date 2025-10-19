package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todo/services/notification-service/internal/email"
	"github.com/todo/services/notification-service/internal/models"
	"github.com/todo/services/notification-service/internal/push"
	"github.com/todo/services/notification-service/internal/repository"
)

type Handler struct {
	repo        *repository.PostgresRepository
	emailSender *email.EmailSender
	pushSender  *push.PushSender
}

func NewHandler(repo *repository.PostgresRepository, emailSender *email.EmailSender, pushSender *push.PushSender) *Handler {
	return &Handler{
		repo:        repo,
		emailSender: emailSender,
		pushSender:  pushSender,
	}
}

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type SendPushRequest struct {
	DeviceToken string `json:"device_token"`
	Title       string `json:"title"`
	Body        string `json:"body"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.emailSender.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		h.repo.SaveNotification(models.TypeEmail, req.To, req.Subject, req.Body, false)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.repo.SaveNotification(models.TypeEmail, req.To, req.Subject, req.Body, true)

	response := Response{
		Success: true,
		Message: "Email sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SendPush(w http.ResponseWriter, r *http.Request) {
	var req SendPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.pushSender.SendPushNotification(req.DeviceToken, req.Title, req.Body)
	if err != nil {
		h.repo.SaveNotification(models.TypePush, req.DeviceToken, req.Title, req.Body, false)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.repo.SaveNotification(models.TypePush, req.DeviceToken, req.Title, req.Body, true)

	response := Response{
		Success: true,
		Message: "Push notification sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/notifications/email", h.SendEmail).Methods("POST")
	router.HandleFunc("/api/notifications/push", h.SendPush).Methods("POST")
}
