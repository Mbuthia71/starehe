package points

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/services"
)

type PointsHandler struct {
	pointsService *services.PointsService
}

func NewPointsHandler(pointsService *services.PointsService) *PointsHandler {
	return &PointsHandler{
		pointsService: pointsService,
	}
}

// GetBalance returns the user's points balance
func (h *PointsHandler) GetBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	balance, err := h.pointsService.GetBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get balance",
		})
	}

	return c.JSON(balance)
}

// GetHistory returns the user's points transaction history
func (h *PointsHandler) GetHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 100 {
		limit = 20
	}

	cursor := c.Query("cursor")

	history, err := h.pointsService.GetHistory(c.Context(), userID, limit, &cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get history",
		})
	}

	return c.JSON(history)
}

// RedeemPoints handles points redemption
func (h *PointsHandler) RedeemPoints(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req models.RedeemPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.pointsService.RedeemPoints(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

// GetRedemption returns details of a specific redemption
func (h *PointsHandler) GetRedemption(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	redemptionID := c.Params("id")

	// This would need to be implemented in the service layer
	// For now, return a placeholder
	return c.JSON(fiber.Map{
		"redemption_id": redemptionID,
		"user_id": userID,
		"status": "requested",
	})
}

// GetRedemptions returns user's redemption history
func (h *PointsHandler) GetRedemptions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 100 {
		limit = 20
	}

	cursor := c.Query("cursor")

	// This would need to be implemented in the service layer
	// For now, return a placeholder
	return c.JSON(fiber.Map{
		"user_id": userID,
		"redemptions": []interface{}{},
		"cursor": cursor,
		"has_more": false,
	})
}

// CreateReferral creates a new referral code for the user
func (h *PointsHandler) CreateReferral(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	response, err := h.pointsService.CreateReferral(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create referral",
		})
	}

	return c.JSON(response)
}

// GetBadges returns user's badges
func (h *PointsHandler) GetBadges(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	badges, err := h.pointsService.GetUserBadges(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get badges",
		})
	}

	return c.JSON(fiber.Map{
		"badges": badges,
	})
}

// GetTier returns user's current tier
func (h *PointsHandler) GetTier(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	tier, err := h.pointsService.GetUserTier(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tier",
		})
	}

	if tier == nil {
		return c.JSON(fiber.Map{
			"tier": "bronze",
			"points_earned_lifetime": 0,
			"current_level_progress": 0,
		})
	}

	return c.JSON(tier)
}

// GetCampaigns returns active points campaigns
func (h *PointsHandler) GetCampaigns(c *fiber.Ctx) error {
	campaigns, err := h.pointsService.GetActiveCampaigns(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get campaigns",
		})
	}

	return c.JSON(fiber.Map{
		"campaigns": campaigns,
	})
}

// ClaimCampaign allows user to claim points from a campaign
func (h *PointsHandler) ClaimCampaign(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	campaignID := c.Params("id")

	err := h.pointsService.AwardCampaignPoints(c.Context(), userID, campaignID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Campaign points awarded successfully",
	})
}

// Admin handlers

// AdminAdjustPoints allows admins to adjust user points
func (h *PointsHandler) AdminAdjustPoints(c *fiber.Ctx) error {
	var req models.AdjustPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	transaction, err := h.pointsService.AdjustPoints(c.Context(), req.UserID, req.Amount, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(transaction)
}

// AdminCreatePartner allows admins to create a new partner
func (h *PointsHandler) AdminCreatePartner(c *fiber.Ctx) error {
	var req models.CreatePartnerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// This would need to be implemented in the service layer
	// For now, return a placeholder
	return c.JSON(fiber.Map{
		"message": "Partner creation not yet implemented",
	})
}

// AdminGetPartners returns all partners
func (h *PointsHandler) AdminGetPartners(c *fiber.Ctx) error {
	// This would need to be implemented in the service layer
	// For now, return a placeholder
	return c.JSON(fiber.Map{
		"partners": []interface{}{},
	})
}
