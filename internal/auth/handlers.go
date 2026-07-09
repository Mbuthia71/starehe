package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/services"
	"starehian-society-platform/pkg/ratelimit"
)

type AuthHandler struct {
	authService *services.AuthService
	rateLimiter *ratelimit.RateLimiter
}

func NewAuthHandler(authService *services.AuthService, rateLimiter *ratelimit.RateLimiter) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		rateLimiter: rateLimiter,
	}
}

// RequestOTP handles OTP request
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var req struct {
		Phone string `json:"phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	// Get IP address for rate limiting
	ip := c.IP()

	// Send OTP
	err := h.authService.RequestOTP(c.Context(), req.Phone, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}

// Signup handles user signup with OTP
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req services.SignupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" || req.FullName == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone, full name, and OTP are required",
		})
	}

	// Get IP address for rate limiting
	ip := c.IP()

	// Signup with OTP
	resp, err := h.authService.SignupWithOTP(c.Context(), &req, ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// Login handles user login with OTP
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone and OTP are required",
		})
	}

	// Get IP address for rate limiting
	ip := c.IP()

	// Login with OTP
	resp, err := h.authService.LoginWithOTP(c.Context(), &req, ip)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// AdminLogin handles admin login with email and password
func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Login with password
	resp, err := h.authService.LoginWithPassword(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	// Refresh token
	tokenPair, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tokenPair)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	// Logout
	err := h.authService.Logout(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// GetMe returns the current user's info
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get user info from auth service (we'd need to add this method)
	// For now, return user ID
	return c.JSON(fiber.Map{
		"user_id": userID,
	})
}

// SetPassword sets a password for a user (admin only)
func (h *AuthHandler) SetPassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req struct {
		TargetUserID string `json:"target_user_id"`
		Password     string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.TargetUserID == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID and password are required",
		})
	}

	// Set password
	err := h.authService.SetPassword(c.Context(), req.TargetUserID, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password set successfully",
	})
}

// Helper function to get client IP
func getClientIP(c *fiber.Ctx) string {
	// Check for forwarded IP (behind proxy)
	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if forwarded := c.Get("X-Real-IP"); forwarded != "" {
		return forwarded
	}
	return c.IP()
}

// Helper function to parse JSON body
func parseJSONBody(c *fiber.Ctx, v interface{}) error {
	body := c.Body()
	if len(body) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Empty request body")
	}
	return json.Unmarshal(body, v)
}
