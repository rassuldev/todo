package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// Create refresh tokens table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			token VARCHAR(500) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) ValidateCredentials(username, password string) (string, string, error) {
	// Connect to user service database to validate credentials
	var userID, hashedPassword string
	err := r.db.QueryRow(
		"SELECT id, password FROM users WHERE username = $1",
		username,
	).Scan(&userID, &hashedPassword)

	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return "", "", err
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	return userID, username, nil
}

func (r *PostgresRepository) StoreRefreshToken(userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt,
	)
	return err
}

func (r *PostgresRepository) GetRefreshToken(token string) (string, error) {
	var userID string
	var expiresAt time.Time

	err := r.db.QueryRow(
		"SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1",
		token,
	).Scan(&userID, &expiresAt)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid refresh token")
	}
	if err != nil {
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("refresh token expired")
	}

	return userID, nil
}

func (r *PostgresRepository) DeleteRefreshToken(token string) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE token = $1", token)
	return err
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
