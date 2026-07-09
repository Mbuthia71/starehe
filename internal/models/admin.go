package models

import "time"

type Report struct {
	ID          string    `json:"id" db:"id"`
	ReporterID  string    `json:"reporter_id" db:"reporter_id"`
	TargetType  string    `json:"target_type" db:"target_type"`
	TargetID    string    `json:"target_id" db:"target_id"`
	Reason      string    `json:"reason" db:"reason"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusReviewed  ReportStatus = "reviewed"
	ReportStatusResolved  ReportStatus = "resolved"
	ReportStatusDismissed ReportStatus = "dismissed"
)

type AdminAction struct {
	ID         string    `json:"id" db:"id"`
	AdminID    string    `json:"admin_id" db:"admin_id"`
	Action     string    `json:"action" db:"action"`
	TargetType *string   `json:"target_type,omitempty" db:"target_type"`
	TargetID   *string   `json:"target_id,omitempty" db:"target_id"`
	Notes      *string   `json:"notes,omitempty" db:"notes"`
	IPAddress  *string   `json:"ip_address,omitempty" db:"ip_address"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type AlumniRoster struct {
	ID         string     `json:"id" db:"id"`
	FullName   string     `json:"full_name" db:"full_name"`
	ClassYear  int        `json:"class_year" db:"class_year"`
	House      *string    `json:"house,omitempty" db:"house"`
	Phone      *string    `json:"phone,omitempty" db:"phone"`
	Email      *string    `json:"email,omitempty" db:"email"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
