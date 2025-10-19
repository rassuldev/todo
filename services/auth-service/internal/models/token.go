package models

import "time"

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
