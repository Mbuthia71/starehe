package services

import (
	"context"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type AdminService struct {
	userRepo    *repository.UserRepository
	adminRepo   *repository.AdminRepository
	logger      *logger.Logger
}

type CreateReportRequest struct {
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Reason     string `json:"reason"`
}

type UpdateReportRequest struct {
	Status models.ReportStatus `json:"status"`
	Notes  string              `json:"notes"`
}

type AlumniRosterEntry struct {
	FullName  string  `json:"full_name"`
	ClassYear int     `json:"class_year"`
	House     *string `json:"house,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Email     *string `json:"email,omitempty"`
}

func NewAdminService(
	userRepo *repository.UserRepository,
	adminRepo *repository.AdminRepository,
	logger *logger.Logger,
) *AdminService {
	return &AdminService{
		userRepo:  userRepo,
		adminRepo: adminRepo,
		logger:    logger,
	}
}

// LogAdminAction logs an admin action to the audit log
func (s *AdminService) LogAdminAction(ctx context.Context, adminID, action string, targetType, targetID, notes, ipAddress string) error {
	adminAction := &models.AdminAction{
		AdminID:    adminID,
		Action:     action,
		TargetType: &targetType,
		TargetID:   &targetID,
		Notes:      &notes,
		IPAddress:  &ipAddress,
	}

	err := s.adminRepo.CreateAdminAction(ctx, adminAction)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
		return fmt.Errorf("failed to log admin action: %w", err)
	}

	s.logger.Infof("Admin action logged: %s by %s", action, adminID)
	return nil
}

// CreateReport creates a new content report
func (s *AdminService) CreateReport(ctx context.Context, reporterID string, req *CreateReportRequest) (*models.Report, error) {
	report := &models.Report{
		ReporterID: reporterID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
		Status:     string(models.ReportStatusPending),
	}

	err := s.adminRepo.CreateReport(ctx, report)
	if err != nil {
		s.logger.Errorf("Failed to create report: %v", err)
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	s.logger.Infof("Report created: %s by %s", report.ID, reporterID)
	return report, nil
}

// GetReport retrieves a report by ID
func (s *AdminService) GetReport(ctx context.Context, reportID string) (*models.Report, error) {
	report, err := s.adminRepo.GetReport(ctx, reportID)
	if err != nil {
		s.logger.Errorf("Failed to get report: %v", err)
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	if report == nil {
		return nil, fmt.Errorf("report not found")
	}

	return report, nil
}

// UpdateReportStatus updates a report's status
func (s *AdminService) UpdateReportStatus(ctx context.Context, adminID, reportID string, req *UpdateReportRequest, ipAddress string) error {
	err := s.adminRepo.UpdateReportStatus(ctx, reportID, req.Status)
	if err != nil {
		s.logger.Errorf("Failed to update report status: %v", err)
		return fmt.Errorf("failed to update report status: %w", err)
	}

	// Log the action
	err = s.LogAdminAction(ctx, adminID, "update_report_status", "report", reportID, req.Notes, ipAddress)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
	}

	s.logger.Infof("Report status updated: %s to %s by %s", reportID, req.Status, adminID)
	return nil
}

// GetReports retrieves reports with optional status filter
func (s *AdminService) GetReports(ctx context.Context, status *models.ReportStatus, limit, offset int) ([]*models.Report, error) {
	reports, err := s.adminRepo.GetReports(ctx, status, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get reports: %v", err)
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}

// GetAdminActions retrieves admin actions from the audit log
func (s *AdminService) GetAdminActions(ctx context.Context, limit, offset int) ([]*models.AdminAction, error) {
	actions, err := s.adminRepo.GetAdminActions(ctx, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get admin actions: %v", err)
		return nil, fmt.Errorf("failed to get admin actions: %w", err)
	}

	return actions, nil
}

// GetAdminActionsByAdmin retrieves actions by a specific admin
func (s *AdminService) GetAdminActionsByAdmin(ctx context.Context, adminID string, limit, offset int) ([]*models.AdminAction, error) {
	actions, err := s.adminRepo.GetAdminActionsByAdmin(ctx, adminID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get admin actions by admin: %v", err)
		return nil, fmt.Errorf("failed to get admin actions by admin: %w", err)
	}

	return actions, nil
}

// AddAlumniRosterEntry adds a single entry to the alumni roster
func (s *AdminService) AddAlumniRosterEntry(ctx context.Context, entry *AlumniRosterEntry) (*models.AlumniRoster, error) {
	roster := &models.AlumniRoster{
		FullName:  entry.FullName,
		ClassYear: entry.ClassYear,
		House:     entry.House,
		Phone:     entry.Phone,
		Email:     entry.Email,
	}

	err := s.adminRepo.CreateAlumniRosterEntry(ctx, roster)
	if err != nil {
		s.logger.Errorf("Failed to create alumni roster entry: %v", err)
		return nil, fmt.Errorf("failed to create alumni roster entry: %w", err)
	}

	s.logger.Infof("Alumni roster entry created: %s", roster.ID)
	return roster, nil
}

// GetAlumniRoster retrieves the alumni roster
func (s *AdminService) GetAlumniRoster(ctx context.Context, limit, offset int) ([]*models.AlumniRoster, error) {
	roster, err := s.adminRepo.GetAlumniRoster(ctx, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get alumni roster: %v", err)
		return nil, fmt.Errorf("failed to get alumni roster: %w", err)
	}

	return roster, nil
}

// VerifyUserAgainstRoster verifies a user against the alumni roster
func (s *AdminService) VerifyUserAgainstRoster(ctx context.Context, adminID, userID string, ipAddress string) error {
	// Match user to roster
	roster, err := s.adminRepo.MatchUserToRoster(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to match user to roster: %v", err)
		return fmt.Errorf("failed to match user to roster: %w", err)
	}
	if roster == nil {
		return fmt.Errorf("no matching roster entry found")
	}

	// Update user status to active
	err = s.userRepo.UpdateStatus(ctx, userID, models.StatusActive)
	if err != nil {
		s.logger.Errorf("Failed to update user status: %v", err)
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Mark roster entry as verified
	err = s.adminRepo.VerifyAlumniRosterEntry(ctx, roster.ID)
	if err != nil {
		s.logger.Errorf("Failed to verify roster entry: %v", err)
		return fmt.Errorf("failed to verify roster entry: %w", err)
	}

	// Log the action
	err = s.LogAdminAction(ctx, adminID, "verify_user", "user", userID, fmt.Sprintf("Verified against roster entry: %s", roster.ID), ipAddress)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
	}

	s.logger.Infof("User verified against roster: %s", userID)
	return nil
}

// SuspendUser suspends a user account
func (s *AdminService) SuspendUser(ctx context.Context, adminID, userID, notes, ipAddress string) error {
	err := s.userRepo.UpdateStatus(ctx, userID, models.StatusSuspended)
	if err != nil {
		s.logger.Errorf("Failed to suspend user: %v", err)
		return fmt.Errorf("failed to suspend user: %w", err)
	}

	// Log the action
	err = s.LogAdminAction(ctx, adminID, "suspend_user", "user", userID, notes, ipAddress)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
	}

	s.logger.Infof("User suspended: %s by %s", userID, adminID)
	return nil
}

// ActivateUser activates a user account
func (s *AdminService) ActivateUser(ctx context.Context, adminID, userID, notes, ipAddress string) error {
	err := s.userRepo.UpdateStatus(ctx, userID, models.StatusActive)
	if err != nil {
		s.logger.Errorf("Failed to activate user: %v", err)
		return fmt.Errorf("failed to activate user: %w", err)
	}

	// Log the action
	err = s.LogAdminAction(ctx, adminID, "activate_user", "user", userID, notes, ipAddress)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
	}

	s.logger.Infof("User activated: %s by %s", userID, adminID)
	return nil
}

// UpdateUserRole updates a user's role
func (s *AdminService) UpdateUserRole(ctx context.Context, adminID, userID string, role models.UserRole, notes, ipAddress string) error {
	err := s.userRepo.UpdateRole(ctx, userID, role)
	if err != nil {
		s.logger.Errorf("Failed to update user role: %v", err)
		return fmt.Errorf("failed to update user role: %w", err)
	}

	// Log the action
	err = s.LogAdminAction(ctx, adminID, "update_user_role", "user", userID, fmt.Sprintf("Role changed to: %s. %s", role, notes), ipAddress)
	if err != nil {
		s.logger.Errorf("Failed to log admin action: %v", err)
	}

	s.logger.Infof("User role updated: %s to %s by %s", userID, role, adminID)
	return nil
}
