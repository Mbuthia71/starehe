package models

import "time"

type Profile struct {
	UserID            string     `json:"user_id" db:"user_id"`
	FullName          string     `json:"full_name" db:"full_name"`
	Bio              *string    `json:"bio,omitempty" db:"bio"`
	AvatarURL        *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CoverURL         *string    `json:"cover_url,omitempty" db:"cover_url"`
	ClassYear       *int       `json:"class_year,omitempty" db:"class_year"`
	House            *string    `json:"house,omitempty" db:"house"`
	Career           *string    `json:"career,omitempty" db:"career"`
	Location         *string    `json:"location,omitempty" db:"location"`
	ProfileVisibility string    `json:"profile_visibility" db:"profile_visibility"`
	ContactVisibility string    `json:"contact_visibility" db:"contact_visibility"`
	CareerVisibility  string    `json:"career_visibility" db:"career_visibility"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type Visibility string

const (
	VisibilityPublic      Visibility = "public"
	VisibilityConnections Visibility = "connections"
	VisibilityPrivate     Visibility = "private"
)
