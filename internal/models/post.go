package models

import "time"

type Post struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Content     *string   `json:"content,omitempty" db:"content"`
	MediaURLs   []string  `json:"media_urls,omitempty" db:"media_urls"`
	Visibility  string    `json:"visibility" db:"visibility"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Comment struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Reaction struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
