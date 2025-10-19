package models

import "time"

type NotificationType string

const (
	TypeEmail NotificationType = "EMAIL"
	TypePush  NotificationType = "PUSH"
)

type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Recipient string           `json:"recipient"`
	Subject   string           `json:"subject,omitempty"`
	Body      string           `json:"body"`
	Sent      bool             `json:"sent"`
	CreatedAt time.Time        `json:"created_at"`
}
