package models

import "time"

type Notification struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"user_id" db:"user_id"`
	Type      string     `json:"type" db:"type"`
	Payload   map[string]interface{} `json:"payload,omitempty" db:"payload"`
	ReadAt    *time.Time `json:"read_at,omitempty" db:"read_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type NotificationType string

const (
	NotificationTypeConnectionRequest  NotificationType = "connection_request"
	NotificationTypeConnectionAccepted NotificationType = "connection_accepted"
	NotificationTypeLike               NotificationType = "like"
	NotificationTypeComment            NotificationType = "comment"
	NotificationTypeMention            NotificationType = "mention"
	NotificationTypeMessage            NotificationType = "message"
	NotificationTypeAnnouncement       NotificationType = "announcement"
)
