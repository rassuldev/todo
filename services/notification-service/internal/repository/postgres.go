package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/todo/services/notification-service/internal/models"
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
		CREATE TABLE IF NOT EXISTS notifications (
			id VARCHAR(36) PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			recipient VARCHAR(255) NOT NULL,
			subject VARCHAR(255),
			body TEXT NOT NULL,
			sent BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) SaveNotification(notifType models.NotificationType, recipient, subject, body string, sent bool) (*models.Notification, error) {
	notification := &models.Notification{
		ID:        uuid.New().String(),
		Type:      notifType,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		Sent:      sent,
		CreatedAt: time.Now(),
	}

	_, err := r.db.Exec(
		"INSERT INTO notifications (id, type, recipient, subject, body, sent, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		notification.ID, notification.Type, notification.Recipient, notification.Subject, notification.Body, notification.Sent, notification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
