package models

import "time"

type Connection struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	ConnectedUserID string    `json:"connected_user_id" db:"connected_user_id"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ConnectionStatus string

const (
	ConnectionStatusPending  ConnectionStatus = "pending"
	ConnectionStatusAccepted ConnectionStatus = "accepted"
)

type Block struct {
	ID        string    `json:"id" db:"id"`
	BlockerID string    `json:"blocker_id" db:"blocker_id"`
	BlockedID string    `json:"blocked_id" db:"blocked_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
