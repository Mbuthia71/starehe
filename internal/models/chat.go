package models

import "time"

type Conversation struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	Name      *string   `json:"name,omitempty" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct"
	ConversationTypeGroup  ConversationType = "group"
)

type ConversationMember struct {
	ID                 string     `json:"id" db:"id"`
	ConversationID     string     `json:"conversation_id" db:"conversation_id"`
	UserID             string     `json:"user_id" db:"user_id"`
	Role               string     `json:"role" db:"role"`
	JoinedAt           time.Time  `json:"joined_at" db:"joined_at"`
	LastReadMessageID  *string    `json:"last_read_message_id,omitempty" db:"last_read_message_id"`
}

type Message struct {
	ID             string    `json:"id" db:"id"`
	ConversationID *string   `json:"conversation_id,omitempty" db:"conversation_id"`
	GroupID        *string   `json:"group_id,omitempty" db:"group_id"`
	RecipientID    *string   `json:"recipient_id,omitempty" db:"recipient_id"`
	SenderID       string    `json:"sender_id" db:"sender_id"`
	SenderName     *string   `json:"sender_name,omitempty"`
	SenderAvatar   *string   `json:"sender_avatar,omitempty"`
	Content        *string   `json:"content,omitempty" db:"content"`
	MediaURL       *string   `json:"media_url,omitempty" db:"media_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
