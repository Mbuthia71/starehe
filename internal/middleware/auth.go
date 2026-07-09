package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/tokens"
	"starehian-society-platform/pkg/logger"
)

type AuthMiddleware struct {
	jwtService *tokens.JWTService
	logger     *logger.Logger
}

func NewAuthMiddleware(jwtService *tokens.JWTService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// Authenticate validates the JWT token and adds user context
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	token := parts[1]

	// Validate token
	claims, err := m.jwtService.ValidateAccessToken(token)
	if err != nil {
		m.logger.Errorf("Failed to validate token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Add user context to Fiber context
	c.Locals("user_id", claims.UserID)
	c.Locals("user_role", claims.Role)

	return c.Next()
}

// RequireRole checks if the user has the required role
func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found in context",
			})
		}

		if userRole.(string) != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAdmin checks if the user has any admin role
func (m *AuthMiddleware) RequireAdmin(c *fiber.Ctx) error {
	userRole := c.Locals("user_role")
	if userRole == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User role not found in context",
		})
	}

	role := userRole.(string)
	if role != "super_admin" && role != "moderator" && role != "support" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	return c.Next()
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) string {
	if userRole, ok := ctx.Value("user_role").(string); ok {
		return userRole
	}
	return ""
}
