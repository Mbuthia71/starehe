package services

import (
	"context"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type ProfileService struct {
	profileRepo      *repository.ProfileRepository
	userRepo         *repository.UserRepository
	logger           *logger.Logger
}

type UpdateProfileRequest struct {
	FullName          *string       `json:"full_name,omitempty"`
	Bio              *string       `json:"bio,omitempty"`
	AvatarURL        *string       `json:"avatar_url,omitempty"`
	CoverURL         *string       `json:"cover_url,omitempty"`
	ClassYear       *int          `json:"class_year,omitempty"`
	House            *models.House `json:"house,omitempty"`
	Career           *string       `json:"career,omitempty"`
	Location         *string       `json:"location,omitempty"`
	ProfileVisibility *string      `json:"profile_visibility,omitempty"`
	ContactVisibility *string      `json:"contact_visibility,omitempty"`
	CareerVisibility  *string      `json:"career_visibility,omitempty"`
}

type SearchProfilesRequest struct {
	SearchTerm string `json:"search_term,omitempty"`
	ClassYear  *int   `json:"class_year,omitempty"`
	House      string `json:"house,omitempty"`
	Location   string `json:"location,omitempty"`
	Career     string `json:"career,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

func NewProfileService(
	profileRepo *repository.ProfileRepository,
	userRepo *repository.UserRepository,
	logger *logger.Logger,
) *ProfileService {
	return &ProfileService{
		profileRepo:  profileRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// GetProfile retrieves a user's profile with privacy checks
func (s *ProfileService) GetProfile(ctx context.Context, viewerID, targetID string) (*models.Profile, error) {
	// TODO: Add authorization check when authzService is decoupled
	// canView, err := s.authzService.CanViewProfile(ctx, viewerID, targetID)
	// if err != nil {
	// 	s.logger.Errorf("Failed to check profile access: %v", err)
	// 	return nil, fmt.Errorf("failed to check access: %w", err)
	// }
	// if !canView {
	// 	return nil, fmt.Errorf("access denied")
	// }

	profile, err := s.profileRepo.GetByID(ctx, targetID)
	if err != nil {
		s.logger.Errorf("Failed to get profile: %v", err)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	// TODO: Add privacy filtering when authzService is decoupled
	// Filter sensitive fields based on visibility
	// if viewerID != targetID {
	// 	// Check contact visibility
	// 	canViewContact, err := s.authzService.CanViewProfileSection(ctx, viewerID, targetID, "contact")
	// 	if err != nil {
	// 		s.logger.Errorf("Failed to check contact visibility: %v", err)
	// 	}
	// 	if !canViewContact {
	// 		// Don't return contact info (phone/email is in user, not profile)
	// 	}
	//
	// 	// Check career visibility
	// 	canViewCareer, err := s.authzService.CanViewProfileSection(ctx, viewerID, targetID, "career")
	// 	if err != nil {
	// 		s.logger.Errorf("Failed to check career visibility: %v", err)
	// 	}
	// 	if !canViewCareer {
	// 		profile.Career = nil
	// 	}
	// }

	return profile, nil
}

// UpdateProfile updates a user's own profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*models.Profile, error) {
	// Get existing profile
	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get profile: %v", err)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	// Update fields if provided
	if req.FullName != nil {
		profile.FullName = *req.FullName
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL
	}
	if req.CoverURL != nil {
		profile.CoverURL = req.CoverURL
	}
	if req.ClassYear != nil {
		profile.ClassYear = req.ClassYear
	}
	if req.House != nil {
		profile.House = req.House
	}
	if req.Career != nil {
		profile.Career = req.Career
	}
	if req.Location != nil {
		profile.Location = req.Location
	}
	if req.ProfileVisibility != nil {
		profile.ProfileVisibility = *req.ProfileVisibility
	}
	if req.ContactVisibility != nil {
		profile.ContactVisibility = *req.ContactVisibility
	}
	if req.CareerVisibility != nil {
		profile.CareerVisibility = *req.CareerVisibility
	}

	// Update profile
	err = s.profileRepo.Update(ctx, profile)
	if err != nil {
		s.logger.Errorf("Failed to update profile: %v", err)
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	s.logger.Infof("Profile updated for user: %s", userID)

	return profile, nil
}

// SearchProfiles searches for alumni with privacy filtering
func (s *ProfileService) SearchProfiles(ctx context.Context, viewerID string, req *SearchProfilesRequest) ([]*models.Profile, error) {
	// Set default pagination
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Search profiles
	profiles, err := s.profileRepo.Search(ctx, req.SearchTerm, req.ClassYear, req.House, req.Location, req.Career, req.Limit, req.Offset)
	if err != nil {
		s.logger.Errorf("Failed to search profiles: %v", err)
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}

	// TODO: Add privacy filtering when authzService is decoupled
	// Filter results based on privacy
	// var filteredProfiles []*models.Profile
	// for _, profile := range profiles {
	// 	canView, err := s.authzService.CanViewProfile(ctx, viewerID, profile.UserID)
	// 	if err != nil {
	// 		s.logger.Errorf("Failed to check profile access for %s: %v", profile.UserID, err)
	// 		continue
	// 	}
	// 	if canView {
	// 		// Apply section-level filtering
	// 		if viewerID != profile.UserID {
	// 			canViewCareer, _ := s.authzService.CanViewProfileSection(ctx, viewerID, profile.UserID, "career")
	// 			if !canViewCareer {
	// 				profile.Career = nil
	// 			}
	// 		}
	// 		filteredProfiles = append(filteredProfiles, profile)
	// 	}
	// }

	return profiles, nil
}

// GetOwnProfile retrieves the user's own profile
func (s *ProfileService) GetOwnProfile(ctx context.Context, userID string) (*models.Profile, error) {
	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get profile: %v", err)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	return profile, nil
}
