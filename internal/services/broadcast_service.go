package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type BroadcastService struct {
	userRepo     *repository.UserRepository
	profileRepo  *repository.ProfileRepository
	notificationRepo *repository.NotificationRepository
	logger       *logger.Logger
}

type BroadcastRequest struct {
	Title       string   `json:"title"`
	Message     string   `json:"message"`
	TargetType  string   `json:"target_type"` // all, class_year, house, location
	TargetValue string   `json:"target_value,omitempty"`
	ScheduledAt *string  `json:"scheduled_at,omitempty"` // ISO timestamp
}

type Broadcast struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	TargetType  string    `json:"target_type"`
	TargetValue string    `json:"target_value,omitempty"`
	SentAt      time.Time `json:"sent_at"`
	SentCount   int       `json:"sent_count"`
}

func NewBroadcastService(
	userRepo *repository.UserRepository,
	profileRepo *repository.ProfileRepository,
	notificationRepo *repository.NotificationRepository,
	logger *logger.Logger,
) *BroadcastService {
	return &BroadcastService{
		userRepo:         userRepo,
		profileRepo:      profileRepo,
		notificationRepo: notificationRepo,
		logger:          logger,
	}
}

// SendBroadcast sends a broadcast message to targeted users
func (s *BroadcastService) SendBroadcast(ctx context.Context, req *BroadcastRequest) (*Broadcast, error) {
	// Get target users based on criteria
	userIDs, err := s.getTargetUsers(ctx, req.TargetType, req.TargetValue)
	if err != nil {
		s.logger.Errorf("Failed to get target users: %v", err)
		return nil, fmt.Errorf("failed to get target users: %w", err)
	}

	if len(userIDs) == 0 {
		return nil, fmt.Errorf("no target users found")
	}

	// Create notifications for all users
	for _, userID := range userIDs {
		notification := &models.Notification{
			ID:      uuid.New().String(),
			UserID:  userID,
			Type:    string(models.NotificationTypeAnnouncement),
			Payload: map[string]interface{}{
				"title":   req.Title,
				"message": req.Message,
			},
		}

		err = s.notificationRepo.CreateNotification(ctx, notification)
		if err != nil {
			s.logger.Errorf("Failed to create notification for user %s: %v", userID, err)
		}
	}

	broadcast := &Broadcast{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Message:     req.Message,
		TargetType:  req.TargetType,
		TargetValue: req.TargetValue,
		SentAt:      time.Now(),
		SentCount:   len(userIDs),
	}

	s.logger.Infof("Broadcast sent to %d users", len(userIDs))
	return broadcast, nil
}

// getTargetUsers retrieves user IDs based on broadcast criteria
func (s *BroadcastService) getTargetUsers(ctx context.Context, targetType, targetValue string) ([]string, error) {
	var userIDs []string

	switch targetType {
	case "all":
		// Get all active users
		users, err := s.userRepo.List(ctx, 10000, 0) // Large limit for all users
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}

	case "class_year":
		// Get users by class year
		profiles, err := s.profileRepo.Search(ctx, "", &targetValue, "", "", "", 10000, 0)
		if err != nil {
			return nil, err
		}
		for _, profile := range profiles {
			userIDs = append(userIDs, profile.UserID)
		}

	case "house":
		// Get users by house
		profiles, err := s.profileRepo.Search(ctx, "", nil, targetValue, "", "", 10000, 0)
		if err != nil {
			return nil, err
		}
		for _, profile := range profiles {
			userIDs = append(userIDs, profile.UserID)
		}

	case "location":
		// Get users by location
		profiles, err := s.profileRepo.Search(ctx, "", nil, "", targetValue, "", 10000, 0)
		if err != nil {
			return nil, err
		}
		for _, profile := range profiles {
			userIDs = append(userIDs, profile.UserID)
		}

	default:
		return nil, fmt.Errorf("invalid target type: %s", targetType)
	}

	return userIDs, nil
}
