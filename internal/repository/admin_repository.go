package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type AdminRepository struct {
	db *database.DB
}

func NewAdminRepository(db *database.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// CreateAdminAction logs an admin action to the audit log
func (r *AdminRepository) CreateAdminAction(ctx context.Context, action *models.AdminAction) error {
	query := `
		INSERT INTO admin_actions (admin_id, action, target_type, target_id, notes, ip_address)
		VALUES (:admin_id, :action, :target_type, :target_id, :notes, :ip_address)
		RETURNING id, created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, action)
	if err != nil {
		return fmt.Errorf("failed to create admin action: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&action.ID, &action.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan admin action: %w", err)
		}
	}
	
	return nil
}

// GetAdminActions retrieves admin actions with pagination
func (r *AdminRepository) GetAdminActions(ctx context.Context, limit, offset int) ([]*models.AdminAction, error) {
	var actions []*models.AdminAction
	query := `
		SELECT id, admin_id, action, target_type, target_id, notes, ip_address, created_at
		FROM admin_actions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	err := r.db.SelectContext(ctx, &actions, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin actions: %w", err)
	}
	
	return actions, nil
}

// GetAdminActionsByAdmin retrieves actions by a specific admin
func (r *AdminRepository) GetAdminActionsByAdmin(ctx context.Context, adminID string, limit, offset int) ([]*models.AdminAction, error) {
	var actions []*models.AdminAction
	query := `
		SELECT id, admin_id, action, target_type, target_id, notes, ip_address, created_at
		FROM admin_actions
		WHERE admin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &actions, query, adminID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin actions by admin: %w", err)
	}
	
	return actions, nil
}

// CreateReport creates a new report
func (r *AdminRepository) CreateReport(ctx context.Context, report *models.Report) error {
	query := `
		INSERT INTO reports (reporter_id, target_type, target_id, reason, status)
		VALUES (:reporter_id, :target_type, :target_id, :reason, :status)
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, report)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&report.ID, &report.CreatedAt, &report.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan report: %w", err)
		}
	}
	
	return nil
}

// GetReport retrieves a report by ID
func (r *AdminRepository) GetReport(ctx context.Context, id string) (*models.Report, error) {
	var report models.Report
	query := `
		SELECT id, reporter_id, target_type, target_id, reason, status, created_at, updated_at
		FROM reports
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &report, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	
	return &report, nil
}

// UpdateReportStatus updates a report's status
func (r *AdminRepository) UpdateReportStatus(ctx context.Context, reportID string, status models.ReportStatus) error {
	query := `UPDATE reports SET status = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, status, reportID)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("report not found")
	}
	
	return nil
}

// GetReports retrieves reports with pagination and filters
func (r *AdminRepository) GetReports(ctx context.Context, status *models.ReportStatus, limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	
	var query string
	var args []interface{}
	
	if status != nil {
		query = `
			SELECT id, reporter_id, target_type, target_id, reason, status, created_at, updated_at
			FROM reports
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{*status, limit, offset}
	} else {
		query = `
			SELECT id, reporter_id, target_type, target_id, reason, status, created_at, updated_at
			FROM reports
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}
	
	err := r.db.SelectContext(ctx, &reports, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}
	
	return reports, nil
}

// CreateAlumniRosterEntry creates a new alumni roster entry
func (r *AdminRepository) CreateAlumniRosterEntry(ctx context.Context, roster *models.AlumniRoster) error {
	query := `
		INSERT INTO alumni_roster (full_name, class_year, house, phone, email)
		VALUES (:full_name, :class_year, :house, :phone, :email)
		RETURNING id, created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, roster)
	if err != nil {
		return fmt.Errorf("failed to create alumni roster entry: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&roster.ID, &roster.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan alumni roster entry: %w", err)
		}
	}
	
	return nil
}

// GetAlumniRoster retrieves alumni roster entries
func (r *AdminRepository) GetAlumniRoster(ctx context.Context, limit, offset int) ([]*models.AlumniRoster, error) {
	var roster []*models.AlumniRoster
	query := `
		SELECT id, full_name, class_year, house, phone, email, verified_at, created_at
		FROM alumni_roster
		ORDER BY class_year DESC, full_name ASC
		LIMIT $1 OFFSET $2
	`
	
	err := r.db.SelectContext(ctx, &roster, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni roster: %w", err)
	}
	
	return roster, nil
}

// VerifyAlumniRosterEntry marks a roster entry as verified
func (r *AdminRepository) VerifyAlumniRosterEntry(ctx context.Context, rosterID string) error {
	query := `UPDATE alumni_roster SET verified_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, rosterID)
	if err != nil {
		return fmt.Errorf("failed to verify alumni roster entry: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("roster entry not found")
	}
	
	return nil
}

// MatchUserToRoster matches a user to the alumni roster
func (r *AdminRepository) MatchUserToRoster(ctx context.Context, userID string) (*models.AlumniRoster, error) {
	var roster models.AlumniRoster
	query := `
		SELECT ar.id, ar.full_name, ar.class_year, ar.house, ar.phone, ar.email, ar.verified_at, ar.created_at
		FROM alumni_roster ar
		INNER JOIN users u ON u.phone = ar.phone OR u.email = ar.email
		WHERE u.id = $1
		LIMIT 1
	`
	
	err := r.db.GetContext(ctx, &roster, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to match user to roster: %w", err)
	}
	
	return &roster, nil
}
