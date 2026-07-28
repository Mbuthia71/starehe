package profiles

import (
	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
	userRepo       *repository.UserRepository
}

func NewProfileHandler(profileService *services.ProfileService, userRepo *repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		userRepo:       userRepo,
	}
}

type ProfileResponse struct {
	*models.Profile
	FileNumber *string `json:"file_number,omitempty"`
}

// GetProfile retrieves a user's profile
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	viewerID := middleware.GetUserID(c.Context())
	targetID := c.Params("id")

	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	profile, err := h.profileService.GetProfile(c.Context(), viewerID, targetID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user to include file_number
	user, err := h.userRepo.GetByID(c.Context(), targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &ProfileResponse{
		Profile:    profile,
		FileNumber: user.FileNumber,
	}

	return c.JSON(response)
}

// GetOwnProfile retrieves the current user's profile
func (h *ProfileHandler) GetOwnProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	profile, err := h.profileService.GetOwnProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user to include file_number
	user, err := h.userRepo.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &ProfileResponse{
		Profile:    profile,
		FileNumber: user.FileNumber,
	}

	return c.JSON(response)
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	profile, err := h.profileService.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(profile)
}

// SearchProfiles searches for alumni
func (h *ProfileHandler) SearchProfiles(c *fiber.Ctx) error {
	viewerID := middleware.GetUserID(c.Context())
	if viewerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.SearchProfilesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	profiles, err := h.profileService.SearchProfiles(c.Context(), viewerID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add file_number to each profile
	var profileResponses []*ProfileResponse
	for _, profile := range profiles {
		user, err := h.userRepo.GetByID(c.Context(), profile.UserID)
		if err != nil {
			continue // Skip if user not found
		}
		profileResponses = append(profileResponses, &ProfileResponse{
			Profile:    profile,
			FileNumber: user.FileNumber,
		})
	}

	return c.JSON(fiber.Map{
		"profiles": profileResponses,
		"limit":    req.Limit,
		"offset":   req.Offset,
	})
}
