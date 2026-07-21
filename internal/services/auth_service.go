package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/tokens"
	"starehian-society-platform/pkg/logger"
	"starehian-society-platform/pkg/redis"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	redis       *redis.Redis
	jwtService  *tokens.JWTService
	otpService  *tokens.OTPService
	logger      *logger.Logger
}

type SignupRequest struct {
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	OTP      string `json:"otp"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type LoginRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type AuthResponse struct {
	User        *models.User         `json:"user"`
	Profile     *models.Profile      `json:"profile,omitempty"`
	TokenPair   *tokens.TokenPair    `json:"tokens"`
	IsNewUser   bool                 `json:"is_new_user"`
}

func NewAuthService(
	userRepo *repository.UserRepository,
	profileRepo *repository.ProfileRepository,
	redis *redis.Redis,
	jwtService *tokens.JWTService,
	otpService *tokens.OTPService,
	logger *logger.Logger,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		redis:       redis,
		jwtService:  jwtService,
		otpService:  otpService,
		logger:      logger,
	}
}

// RequestOTP sends an OTP to the given phone number
func (s *AuthService) RequestOTP(ctx context.Context, phone, ip string) error {
	// Check if user exists
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		s.logger.Errorf("Failed to check user existence: %v", err)
		return fmt.Errorf("failed to check user: %w", err)
	}

	// Send OTP
	err = s.otpService.SendOTP(ctx, phone, ip)
	if err != nil {
		return err
	}

	if user != nil {
		s.logger.Infof("OTP requested for existing user: %s", phone)
	} else {
		s.logger.Infof("OTP requested for new user: %s", phone)
	}

	return nil
}

// SignupWithPassword creates a new user with password
func (s *AuthService) SignupWithPassword(ctx context.Context, phone, fullName, password, ip string) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	hashedPasswordStr := string(hashedPassword)
	user := &models.User{
		ID:           uuid.New().String(),
		Phone:        phone,
		PasswordHash: &hashedPasswordStr,
		Role:         string(models.RoleMember),
		Status:       string(models.StatusActive), // Auto-activate for password signup
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create profile
	profile := &models.Profile{
		UserID:            user.ID,
		FullName:          fullName,
		ProfileVisibility: string(models.VisibilityConnections),
		ContactVisibility: string(models.VisibilityConnections),
		CareerVisibility:  string(models.VisibilityConnections),
	}

	err = s.profileRepo.Create(ctx, profile)
	if err != nil {
		s.logger.Errorf("Failed to create profile: %v", err)
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens: %v", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store session in Redis
	err = s.redis.SetSession(ctx, tokenPair.RefreshToken, user.ID, 7*24*time.Hour)
	if err != nil {
		s.logger.Errorf("Failed to store session: %v", err)
	}

	s.logger.Infof("New user created: %s", user.ID)

	return &AuthResponse{
		User:      user,
		Profile:   profile,
		TokenPair: tokenPair,
		IsNewUser: true,
	}, nil
}

// LoginWithPassword authenticates a user with password
func (s *AuthService) LoginWithPassword(ctx context.Context, phone, password, ip string) (*AuthResponse, error) {
	// Get user
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user has password
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, fmt.Errorf("user has no password set. Please use OTP or contact admin")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Check user status
	if user.Status != string(models.StatusActive) {
		return nil, fmt.Errorf("user account is not active")
	}

	// Get profile
	profile, err := s.profileRepo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorf("Failed to get profile: %v", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens: %v", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store session in Redis
	err = s.redis.SetSession(ctx, tokenPair.RefreshToken, user.ID, 7*24*time.Hour)
	if err != nil {
		s.logger.Errorf("Failed to store session: %v", err)
	}

	s.logger.Infof("User logged in: %s", user.ID)

	return &AuthResponse{
		User:      user,
		Profile:   profile,
		TokenPair: tokenPair,
		IsNewUser: false,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*tokens.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if session exists in Redis
	userID, err := s.redis.GetSession(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	if userID != claims.UserID {
		return nil, fmt.Errorf("session mismatch")
	}

	// Get user to get current role
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens: %v", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete old session and create new one
	s.redis.DeleteSession(ctx, refreshToken)
	err = s.redis.SetSession(ctx, tokenPair.RefreshToken, user.ID, 7*24*time.Hour)
	if err != nil {
		s.logger.Errorf("Failed to store session: %v", err)
	}

	s.logger.Infof("Token refreshed for user: %s", user.ID)

	return tokenPair, nil
}

// Logout logs out a user by deleting their session
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	err := s.redis.DeleteSession(ctx, refreshToken)
	if err != nil {
		s.logger.Errorf("Failed to delete session: %v", err)
		return fmt.Errorf("failed to logout: %w", err)
	}

	s.logger.Infof("User logged out")
	return nil
}

// SetPassword sets a password for a user (for admin login)
func (s *AuthService) SetPassword(ctx context.Context, userID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	hashedStr := string(hashedPassword)
	err = s.userRepo.Update(ctx, &models.User{
		ID:           userID,
		PasswordHash: &hashedStr,
	})
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

// LoginWithEmail authenticates a user with email and password (for admin)
func (s *AuthService) LoginWithEmail(ctx context.Context, email, password string) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if user.PasswordHash == nil {
		return nil, fmt.Errorf("password not set for this account")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != string(models.StatusActive) {
		return nil, fmt.Errorf("user account is not active")
	}

	// Get profile
	profile, err := s.profileRepo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorf("Failed to get profile: %v", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens: %v", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store session in Redis
	err = s.redis.SetSession(ctx, tokenPair.RefreshToken, user.ID, 7*24*time.Hour)
	if err != nil {
		s.logger.Errorf("Failed to store session: %v", err)
	}

	s.logger.Infof("Admin logged in: %s", user.ID)

	return &AuthResponse{
		User:      user,
		Profile:   profile,
		TokenPair: tokenPair,
		IsNewUser: false,
	}, nil
}
