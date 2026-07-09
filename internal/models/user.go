package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Phone        string    `json:"phone" db:"phone"`
	Email        *string   `json:"email,omitempty" db:"email"`
	PasswordHash *string   `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleModerator  UserRole = "moderator"
	RoleSupport    UserRole = "support"
	RoleMember     UserRole = "member"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusDeleted  UserStatus = "deleted"
	StatusPending  UserStatus = "pending"
)
