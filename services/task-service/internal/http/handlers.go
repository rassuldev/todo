package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/todo/services/task-service/internal/models"
	"github.com/todo/services/task-service/internal/repository"
)

type Handler struct {
	repo *repository.PostgresRepository
}

func NewHandler(repo *repository.PostgresRepository) *Handler {
	return &Handler{repo: repo}
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	UserID      string  `json:"user_id"`
	DueDate     *string `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	DueDate     *string `json:"due_date,omitempty"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	priority := models.TaskPriority(req.Priority)
	var dueDate *time.Time
	if req.DueDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err == nil {
			dueDate = &parsedDate
		}
	}

	task, err := h.repo.CreateTask(req.Title, req.Description, req.UserID, priority, dueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := models.TaskStatus(req.Status)
	priority := models.TaskPriority(req.Priority)
	var dueDate *time.Time
	if req.DueDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err == nil {
			dueDate = &parsedDate
		}
	}

	task, err := h.repo.UpdateTask(id, req.Title, req.Description, status, priority, dueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.repo.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}

	tasks, total, err := h.repo.ListTasks(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
		"total": total,
		"page":  page,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ListUserTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	statusStr := r.URL.Query().Get("status")

	page := 1
	pageSize := 10
	var status models.TaskStatus
	if statusStr != "" {
		status = models.TaskStatus(statusStr)
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}

	tasks, total, err := h.repo.ListUserTasks(userID, page, pageSize, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
		"total": total,
		"page":  page,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/tasks", h.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks", h.ListTasks).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", h.GetTask).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", h.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}", h.DeleteTask).Methods("DELETE")
	router.HandleFunc("/api/users/{user_id}/tasks", h.ListUserTasks).Methods("GET")
}
