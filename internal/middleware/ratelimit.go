package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/pkg/ratelimit"
)

const (
	// AuthRateLimit: 10 requests per minute per IP
	AuthRateLimit = "auth"
	// GeneralRateLimit: 100 requests per minute per user
	GeneralRateLimit = "general"
	// UploadRateLimit: 5 requests per minute per user
	UploadRateLimit = "upload"
	// AdminRateLimit: 50 requests per minute per admin user
	AdminRateLimit = "admin"
	// ChatRateLimit: 200 requests per minute per user (lenient for real-time messaging)
	ChatRateLimit = "chat"
)

var RateLimitConfig = map[string]int{
	AuthRateLimit:    10,
	GeneralRateLimit: 100,
	UploadRateLimit:  5,
	AdminRateLimit:   50,
	ChatRateLimit:    200, // More lenient for real-time chat
}

func RateLimitMiddleware(limiter *ratelimit.RateLimiter, limitType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("%s:%s", limitType, c.IP())
		userID := GetUserID(c.Context())
		if userID != "" {
			key = fmt.Sprintf("%s:%s", limitType, userID)
		}

		limit, exists := RateLimitConfig[limitType]
		if !exists {
			limit = 100 // Default to general rate limit if not found
		}

		allowed, _, err := limiter.CheckRateLimit(c.Context(), key, limit, time.Minute)
		if err != nil {
			// Fail open if the limiter backend is unavailable
			return c.Next()
		}
		if !allowed {
			return c.Status(429).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		}

		return c.Next()
	}
}
