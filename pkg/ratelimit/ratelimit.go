package ratelimit

import (
	"context"
	"fmt"
	"time"

	"starehian-society-platform/pkg/redis"
)

type RateLimiter struct {
	redis *redis.Redis
}

func NewRateLimiter(redis *redis.Redis) *RateLimiter {
	return &RateLimiter{redis: redis}
}

// CheckRateLimit checks if the request is within rate limits
// Returns true if allowed, false if rate limit exceeded
func (rl *RateLimiter) CheckRateLimit(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, int64, error) {
	count, err := rl.redis.IncrementRateLimit(ctx, key, window)
	if err != nil {
		return false, 0, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	allowed := count <= int64(maxRequests)
	return allowed, count, nil
}

// OTP-specific rate limiting
func (rl *RateLimiter) CheckOTPRateLimit(ctx context.Context, phone string) (bool, error) {
	// Max 3 OTP requests per hour per phone
	allowed, _, err := rl.CheckRateLimit(ctx, "otp:phone:"+phone, 3, time.Hour)
	if err != nil {
		return false, err
	}
	if !allowed {
		return false, nil
	}

	// Max 10 OTP requests per day per phone
	allowed, _, err = rl.CheckRateLimit(ctx, "otp:phone:daily:"+phone, 10, 24*time.Hour)
	if err != nil {
		return false, err
	}

	return allowed, nil
}

// IP-based rate limiting for OTP
func (rl *RateLimiter) CheckIPRateLimit(ctx context.Context, ip string) (bool, error) {
	// Max 5 OTP requests per hour per IP
	allowed, _, err := rl.CheckRateLimit(ctx, "otp:ip:"+ip, 5, time.Hour)
	if err != nil {
		return false, err
	}
	if !allowed {
		return false, nil
	}

	// Max 20 OTP requests per day per IP
	allowed, _, err = rl.CheckRateLimit(ctx, "otp:ip:daily:"+ip, 20, 24*time.Hour)
	return allowed, err
}

// General API rate limiting
func (rl *RateLimiter) CheckAPIRateLimit(ctx context.Context, userID string, endpoint string) (bool, error) {
	// Max 100 requests per minute per user per endpoint
	allowed, _, err := rl.CheckRateLimit(ctx, fmt.Sprintf("api:%s:%s", userID, endpoint), 100, time.Minute)
	return allowed, err
}

// Global rate limiting (DDoS protection)
func (rl *RateLimiter) CheckGlobalRateLimit(ctx context.Context, endpoint string) (bool, error) {
	// Max 1000 requests per second per endpoint globally
	allowed, _, err := rl.CheckRateLimit(ctx, "global:"+endpoint, 1000, time.Second)
	return allowed, err
}
