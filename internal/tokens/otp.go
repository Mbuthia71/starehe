package tokens

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"starehian-society-platform/pkg/logger"
	"starehian-society-platform/pkg/ratelimit"
	"starehian-society-platform/pkg/redis"
)

type OTPService struct {
	redis      *redis.Redis
	rateLimiter *ratelimit.RateLimiter
	atService  *AfricasTalkingService
	logger     *logger.Logger
}

func NewOTPService(redis *redis.Redis, rateLimiter *ratelimit.RateLimiter, atService *AfricasTalkingService, logger *logger.Logger) *OTPService {
	return &OTPService{
		redis:       redis,
		rateLimiter: rateLimiter,
		atService:   atService,
		logger:      logger,
	}
}

// GenerateOTP generates a 6-digit OTP code
func (s *OTPService) GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}
		otp += n.String()
	}
	return otp, nil
}

// SendOTP sends an OTP to the given phone number with rate limiting
func (s *OTPService) SendOTP(ctx context.Context, phone, ip string) error {
	// Format phone number
	formattedPhone := FormatPhone(phone)
	
	// Check rate limits
	allowed, err := s.rateLimiter.CheckOTPRateLimit(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to check OTP rate limit for phone %s: %v", formattedPhone, err)
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	if !allowed {
		s.logger.Infof("OTP rate limit exceeded for phone %s", formattedPhone)
		return fmt.Errorf("too many OTP requests for this phone number")
	}
	
	allowed, err = s.rateLimiter.CheckIPRateLimit(ctx, ip)
	if err != nil {
		s.logger.Errorf("Failed to check IP rate limit: %v", err)
		return fmt.Errorf("failed to check IP rate limit: %w", err)
	}
	if !allowed {
		s.logger.Infof("OTP rate limit exceeded for IP %s", ip)
		return fmt.Errorf("too many OTP requests from this IP")
	}
	
	// Generate OTP
	otp, err := s.GenerateOTP()
	if err != nil {
		s.logger.Errorf("Failed to generate OTP: %v", err)
		return fmt.Errorf("failed to generate OTP: %w", err)
	}
	
	// Store OTP in Redis with 10 minute expiry
	err = s.redis.SetOTP(ctx, formattedPhone, otp, 10*time.Minute)
	if err != nil {
		s.logger.Errorf("Failed to store OTP in Redis: %v", err)
		return fmt.Errorf("failed to store OTP: %w", err)
	}
	
	// Send OTP via Africa's Talking
	err = s.atService.SendOTP(formattedPhone, otp)
	if err != nil {
		s.logger.Errorf("Failed to send OTP via Africa's Talking: %v", err)
		// Delete OTP from Redis if send failed
		s.redis.DeleteOTP(ctx, formattedPhone)
		return fmt.Errorf("failed to send OTP: %w", err)
	}
	
	s.logger.Infof("OTP sent successfully to %s", formattedPhone)
	return nil
}

// VerifyOTP verifies an OTP code for a phone number
func (s *OTPService) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
	formattedPhone := FormatPhone(phone)
	
	// Get stored OTP
	storedOTP, err := s.redis.GetOTP(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to get OTP from Redis: %v", err)
		return false, fmt.Errorf("failed to retrieve OTP: %w", err)
	}
	
	// Verify OTP
	if storedOTP != otp {
		s.logger.Infof("Invalid OTP attempt for phone %s", formattedPhone)
		return false, nil
	}
	
	// Delete OTP after successful verification
	err = s.redis.DeleteOTP(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to delete OTP from Redis: %v", err)
		// Don't return error here as verification was successful
	}
	
	s.logger.Infof("OTP verified successfully for phone %s", formattedPhone)
	return true, nil
}
