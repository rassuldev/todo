package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/todo/services/task-service/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
			priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
			user_id VARCHAR(36) NOT NULL,
			due_date TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) CreateTask(title, description, userID string, priority models.TaskPriority, dueDate *time.Time) (*models.Task, error) {
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      models.StatusPending,
		Priority:    priority,
		UserID:      userID,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := r.db.Exec(
		"INSERT INTO tasks (id, title, description, status, priority, user_id, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		task.ID, task.Title, task.Description, task.Status, task.Priority, task.UserID, task.DueDate, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *PostgresRepository) GetTaskByID(id string) (*models.Task, error) {
	task := &models.Task{}
	var dueDate sql.NullTime

	err := r.db.QueryRow(
		"SELECT id, title, description, status, priority, user_id, due_date, created_at, updated_at FROM tasks WHERE id = $1",
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.UserID, &dueDate, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}

	return task, nil
}

func (r *PostgresRepository) UpdateTask(id, title, description string, status models.TaskStatus, priority models.TaskPriority, dueDate *time.Time) (*models.Task, error) {
	_, err := r.db.Exec(
		"UPDATE tasks SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = $7 WHERE id = $1",
		id, title, description, status, priority, dueDate, time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetTaskByID(id)
}

func (r *PostgresRepository) DeleteTask(id string) error {
	result, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *PostgresRepository) ListTasks(page, pageSize int) ([]*models.Task, int, error) {
	offset := (page - 1) * pageSize

	rows, err := r.db.Query(
		"SELECT id, title, description, status, priority, user_id, due_date, created_at, updated_at FROM tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		var dueDate sql.NullTime

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.UserID, &dueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, task)
	}

	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *PostgresRepository) ListUserTasks(userID string, page, pageSize int, status models.TaskStatus) ([]*models.Task, int, error) {
	offset := (page - 1) * pageSize

	query := "SELECT id, title, description, status, priority, user_id, due_date, created_at, updated_at FROM tasks WHERE user_id = $1"
	countQuery := "SELECT COUNT(*) FROM tasks WHERE user_id = $1"
	args := []interface{}{userID}

	if status != "" {
		query += " AND status = $2"
		countQuery += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		var dueDate sql.NullTime

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.UserID, &dueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, task)
	}

	var total int
	countArgs := args[:len(args)-2] // Remove limit and offset
	err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
