package admin

import (
	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/services"
)

// AnalyticsHandler handles analytics endpoints
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboardMetrics retrieves dashboard metrics
func (h *AnalyticsHandler) GetDashboardMetrics(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	metrics, err := h.analyticsService.GetDashboardMetrics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(metrics)
}

// GetCohortAnalytics retrieves analytics by class year
func (h *AnalyticsHandler) GetCohortAnalytics(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	cohorts, err := h.analyticsService.GetCohortAnalytics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"cohorts": cohorts,
	})
}

// GetSignupsOverTime retrieves signup trends
func (h *AnalyticsHandler) GetSignupsOverTime(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	days := c.QueryInt("days", 30)
	if days > 365 {
		days = 365
	}

	data, err := h.analyticsService.GetSignupsOverTime(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// GetEngagementOverTime retrieves engagement trends
func (h *AnalyticsHandler) GetEngagementOverTime(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	days := c.QueryInt("days", 30)
	if days > 365 {
		days = 365
	}

	data, err := h.analyticsService.GetEngagementOverTime(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// GetTopContent retrieves most popular posts
func (h *AnalyticsHandler) GetTopContent(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	content, err := h.analyticsService.GetTopContent(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"content": content,
	})
}

// BroadcastHandler handles broadcast endpoints
type BroadcastHandler struct {
	broadcastService *services.BroadcastService
}

func NewBroadcastHandler(broadcastService *services.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{
		broadcastService: broadcastService,
	}
}

// SendBroadcast sends a broadcast message
func (h *BroadcastHandler) SendBroadcast(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.BroadcastRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and message are required",
		})
	}

	broadcast, err := h.broadcastService.SendBroadcast(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(broadcast)
}
