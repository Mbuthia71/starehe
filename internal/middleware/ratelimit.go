package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/pkg/ratelimit"
)

type RateLimitConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
	BurstSize         int
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(rateLimiter *ratelimit.RateLimiter, config *RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		userID := c.Locals("user_id")
		
		// Use user ID if available, otherwise use IP
		key := fmt.Sprintf("ratelimit:%s", userID)
		if userID == nil {
			key = fmt.Sprintf("ratelimit:ip:%s", ip)
		}

		// Check minute limit
		allowed, err := rateLimiter.CheckRateLimit(c.Context(), key, config.RequestsPerMinute, time.Minute)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limit check failed",
			})
		}
		if !allowed {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded (per minute)",
			})
		}

		// Check hour limit
		hourKey := key + ":hour"
		allowed, err = rateLimiter.CheckRateLimit(c.Context(), hourKey, config.RequestsPerHour, time.Hour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rate limit check failed",
			})
		}
		if !allowed {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded (per hour)",
			})
		}

		return c.Next()
	}
}

// Rate limit configurations for different endpoint types
var (
	AuthRateLimit = &RateLimitConfig{
		RequestsPerMinute: 10,
		RequestsPerHour:   100,
		BurstSize:         5,
	}

	GeneralRateLimit = &RateLimitConfig{
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
		BurstSize:         20,
	}

	UploadRateLimit = &RateLimitConfig{
		RequestsPerMinute: 10,
		RequestsPerHour:   100,
		BurstSize:         5,
	}

	AdminRateLimit = &RateLimitConfig{
		RequestsPerMinute: 30,
		RequestsPerHour:   500,
		BurstSize:         10,
	}
)
