package admin

import (
	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/services"
)

type AdminHandler struct {
	adminService  *services.AdminService
	userRepo      *repository.UserRepository
	authzService  *middleware.AuthorizationService
}

func NewAdminHandler(
	adminService *services.AdminService,
	userRepo *repository.UserRepository,
	authzService *middleware.AuthorizationService,
) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		userRepo:     userRepo,
		authzService: authzService,
	}
}

// GetUsers retrieves a list of users (admin only)
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	users, err := h.userRepo.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser retrieves a specific user (admin only)
func (h *AdminHandler) GetUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	user, err := h.userRepo.GetByID(c.Context(), targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// SuspendUser suspends a user account (admin only)
func (h *AdminHandler) SuspendUser(c *fiber.Ctx) error {
	adminID := middleware.GetUserID(c.Context())
	if adminID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), adminID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	ip := c.IP()
	err = h.adminService.SuspendUser(c.Context(), adminID, targetID, req.Notes, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User suspended successfully",
	})
}

// ActivateUser activates a user account (admin only)
func (h *AdminHandler) ActivateUser(c *fiber.Ctx) error {
	adminID := middleware.GetUserID(c.Context())
	if adminID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), adminID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	ip := c.IP()
	err = h.adminService.ActivateUser(c.Context(), adminID, targetID, req.Notes, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User activated successfully",
	})
}

// UpdateUserRole updates a user's role (super admin only)
func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	adminID := middleware.GetUserID(c.Context())
	if adminID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is super admin
	isSuperAdmin, err := h.authzService.IsSuperAdmin(c.Context(), adminID)
	if err != nil || !isSuperAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Super admin access required",
		})
	}

	var req struct {
		Role  models.UserRole `json:"role"`
		Notes string          `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	ip := c.IP()
	err = h.adminService.UpdateUserRole(c.Context(), adminID, targetID, req.Role, req.Notes, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User role updated successfully",
	})
}

// CreateReport creates a content report
func (h *AdminHandler) CreateReport(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	report, err := h.adminService.CreateReport(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// GetReports retrieves reports (admin only)
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	statusStr := c.Query("status", "")

	var status *models.ReportStatus
	if statusStr != "" {
		s := models.ReportStatus(statusStr)
		status = &s
	}

	reports, err := h.adminService.GetReports(c.Context(), status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"reports": reports,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateReportStatus updates a report's status (admin only)
func (h *AdminHandler) UpdateReportStatus(c *fiber.Ctx) error {
	adminID := middleware.GetUserID(c.Context())
	if adminID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), adminID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	reportID := c.Params("id")
	if reportID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Report ID is required",
		})
	}

	var req services.UpdateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ip := c.IP()
	err = h.adminService.UpdateReportStatus(c.Context(), adminID, reportID, &req, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report status updated successfully",
	})
}

// GetAdminActions retrieves admin actions from audit log (admin only)
func (h *AdminHandler) GetAdminActions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	actions, err := h.adminService.GetAdminActions(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"actions": actions,
		"limit":   limit,
		"offset":  offset,
	})
}

// AddAlumniRosterEntry adds a single entry to the alumni roster (admin only)
func (h *AdminHandler) AddAlumniRosterEntry(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req services.AlumniRosterEntry
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	roster, err := h.adminService.AddAlumniRosterEntry(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(roster)
}

// GetAlumniRoster retrieves the alumni roster (admin only)
func (h *AdminHandler) GetAlumniRoster(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	roster, err := h.adminService.GetAlumniRoster(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"roster": roster,
		"limit":  limit,
		"offset": offset,
	})
}

// VerifyUserAgainstRoster verifies a user against the alumni roster (admin only)
func (h *AdminHandler) VerifyUserAgainstRoster(c *fiber.Ctx) error {
	adminID := middleware.GetUserID(c.Context())
	if adminID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user is admin
	isAdmin, err := h.authzService.IsAdmin(c.Context(), adminID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	ip := c.IP()
	err = h.adminService.VerifyUserAgainstRoster(c.Context(), adminID, targetID, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User verified successfully",
	})
}
