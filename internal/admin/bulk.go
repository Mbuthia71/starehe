package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type BulkOperationService struct {
	userRepo    *repository.UserRepository
	adminRepo   *repository.AdminRepository
	logger      *logger.Logger
}

type BulkSuspendRequest struct {
	UserIDs []string `json:"user_ids"`
	Notes   string   `json:"notes"`
}

type BulkActivateRequest struct {
	UserIDs []string `json:"user_ids"`
	Notes   string   `json:"notes"`
}

type BulkVerifyRequest struct {
	UserIDs []string `json:"user_ids"`
}

type BulkDeleteRequest struct {
	UserIDs []string `json:"user_ids"`
	Notes   string   `json:"notes"`
}

type BulkOperationResult struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedIDs    []string `json:"failed_ids,omitempty"`
}

func NewBulkOperationService(
	userRepo *repository.UserRepository,
	adminRepo *repository.AdminRepository,
	logger *logger.Logger,
) *BulkOperationService {
	return &BulkOperationService{
		userRepo:  userRepo,
		adminRepo: adminRepo,
		logger:    logger,
	}
}

// BulkSuspend suspends multiple users at once
func (s *BulkOperationService) BulkSuspend(ctx context.Context, adminID string, req *BulkSuspendRequest, ipAddress string) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		FailedIDs: []string{},
	}

	for _, userID := range req.UserIDs {
		err := s.userRepo.UpdateStatus(ctx, userID, models.StatusSuspended)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Errorf("Failed to suspend user %s: %v", userID, err)
			continue
		}

		// Log action
		s.adminRepo.CreateAdminAction(ctx, &models.AdminAction{
			ID:         uuid.New().String(),
			AdminID:    adminID,
			Action:     "bulk_suspend_user",
			TargetType: &[]string{"user"}[0],
			TargetID:   &userID,
			Notes:      &req.Notes,
			IPAddress:  &ipAddress,
		})

		result.SuccessCount++
	}

	s.logger.Infof("Bulk suspend completed: %d succeeded, %d failed", result.SuccessCount, result.FailedCount)
	return result, nil
}

// BulkActivate activates multiple users at once
func (s *BulkOperationService) BulkActivate(ctx context.Context, adminID string, req *BulkActivateRequest, ipAddress string) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		FailedIDs: []string{},
	}

	for _, userID := range req.UserIDs {
		err := s.userRepo.UpdateStatus(ctx, userID, models.StatusActive)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Errorf("Failed to activate user %s: %v", userID, err)
			continue
		}

		// Log action
		s.adminRepo.CreateAdminAction(ctx, &models.AdminAction{
			ID:         uuid.New().String(),
			AdminID:    adminID,
			Action:     "bulk_activate_user",
			TargetType: &[]string{"user"}[0],
			TargetID:   &userID,
			Notes:      &req.Notes,
			IPAddress:  &ipAddress,
		})

		result.SuccessCount++
	}

	s.logger.Infof("Bulk activate completed: %d succeeded, %d failed", result.SuccessCount, result.FailedCount)
	return result, nil
}

// BulkVerify verifies multiple users against the alumni roster
func (s *BulkOperationService) BulkVerify(ctx context.Context, adminID string, req *BulkVerifyRequest, ipAddress string) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		FailedIDs: []string{},
	}

	for _, userID := range req.UserIDs {
		// Match user to roster
		roster, err := s.adminRepo.MatchUserToRoster(ctx, userID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Errorf("Failed to match user %s to roster: %v", userID, err)
			continue
		}
		if roster == nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Infof("No roster match for user %s", userID)
			continue
		}

		// Update user status
		err = s.userRepo.UpdateStatus(ctx, userID, models.StatusActive)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Errorf("Failed to activate user %s: %v", userID, err)
			continue
		}

		// Mark roster entry as verified
		err = s.adminRepo.VerifyAlumniRosterEntry(ctx, roster.ID)
		if err != nil {
			s.logger.Errorf("Failed to verify roster entry %s: %v", roster.ID, err)
		}

		// Log action
		s.adminRepo.CreateAdminAction(ctx, &models.AdminAction{
			ID:         uuid.New().String(),
			AdminID:    adminID,
			Action:     "bulk_verify_user",
			TargetType: &[]string{"user"}[0],
			TargetID:   &userID,
			Notes:      &[]string{fmt.Sprintf("Verified against roster entry: %s", roster.ID)}[0],
			IPAddress:  &ipAddress,
		})

		result.SuccessCount++
	}

	s.logger.Infof("Bulk verify completed: %d succeeded, %d failed", result.SuccessCount, result.FailedCount)
	return result, nil
}

// BulkDelete soft deletes multiple users
func (s *BulkOperationService) BulkDelete(ctx context.Context, adminID string, req *BulkDeleteRequest, ipAddress string) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		FailedIDs: []string{},
	}

	for _, userID := range req.UserIDs {
		err := s.userRepo.UpdateStatus(ctx, userID, models.StatusDeleted)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			s.logger.Errorf("Failed to delete user %s: %v", userID, err)
			continue
		}

		// Log action
		s.adminRepo.CreateAdminAction(ctx, &models.AdminAction{
			ID:         uuid.New().String(),
			AdminID:    adminID,
			Action:     "bulk_delete_user",
			TargetType: &[]string{"user"}[0],
			TargetID:   &userID,
			Notes:      &req.Notes,
			IPAddress:  &ipAddress,
		})

		result.SuccessCount++
	}

	s.logger.Infof("Bulk delete completed: %d succeeded, %d failed", result.SuccessCount, result.FailedCount)
	return result, nil
}
