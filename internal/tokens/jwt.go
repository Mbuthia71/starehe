package tokens

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"starehian-society-platform/pkg/config"
	"starehian-society-platform/pkg/logger"
)

type JWTService struct {
	secret        string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	redisClient   *redis.Client
	logger        *logger.Logger
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret:        cfg.Secret,
		refreshSecret: cfg.RefreshSecret,
		accessExpiry:  time.Duration(cfg.AccessTokenExpiry) * time.Hour,
		refreshExpiry: time.Duration(cfg.RefreshTokenExpiry) * 24 * time.Hour,
	}
}

func NewJWTServiceWithRedis(cfg *config.JWTConfig, redisClient *redis.Client, appLogger *logger.Logger) *JWTService {
	return &JWTService{
		secret:        cfg.Secret,
		refreshSecret: cfg.RefreshSecret,
		accessExpiry:  time.Duration(cfg.AccessTokenExpiry) * time.Hour,
		refreshExpiry: time.Duration(cfg.RefreshTokenExpiry) * 24 * time.Hour,
		redisClient:   redisClient,
		logger:        appLogger,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(userID, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
	}, nil
}

// generateAccessToken generates a JWT access token
func (s *JWTService) generateAccessToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// generateRefreshToken generates a JWT refresh token
func (s *JWTService) generateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

// ValidateAccessToken validates an access token and returns the claims
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.secret)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.refreshSecret)
}

// validateToken validates a JWT token with the given secret
func (s *JWTService) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if user is revoked
		if s.redisClient != nil {
			ctx := context.Background()
			revoked, _ := s.IsUserRevoked(ctx, claims.UserID)
			if revoked {
				return nil, errors.New("user tokens have been revoked")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *JWTService) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token is blacklisted
	if s.redisClient != nil {
		ctx := context.Background()
		tokenHash := s.hashToken(refreshTokenString)
		blacklisted, _ := s.redisClient.Exists(ctx, fmt.Sprintf("blacklist:refresh:%s", tokenHash)).Result()
		if blacklisted > 0 {
			return nil, errors.New("refresh token has been revoked")
		}
	}

	// Generate new token pair
	return s.GenerateTokenPair(claims.UserID, claims.Role)
}

// RevokeToken revokes a token by adding it to the Redis blacklist
func (s *JWTService) RevokeToken(ctx context.Context, tokenString string, tokenType string) error {
	if s.redisClient == nil {
		s.logger.Warn("Redis client not configured, token revocation skipped")
		return nil
	}

	tokenHash := s.hashToken(tokenString)
	var ttl time.Duration

	if tokenType == "access" {
		ttl = s.accessExpiry
	} else if tokenType == "refresh" {
		ttl = s.refreshExpiry
	} else {
		return errors.New("invalid token type")
	}

	key := fmt.Sprintf("blacklist:%s:%s", tokenType, tokenHash)
	err := s.redisClient.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		s.logger.Errorf("Failed to blacklist token: %v", err)
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	s.logger.Infof("Token %s blacklisted", tokenType)
	return nil
}

// RevokeUserTokens revokes all tokens for a specific user
func (s *JWTService) RevokeUserTokens(ctx context.Context, userID string) error {
	if s.redisClient == nil {
		s.logger.Warn("Redis client not configured, user token revocation skipped")
		return nil
	}

	// Add user to revoked users set with TTL
	key := fmt.Sprintf("revoked_users:%s", userID)
	err := s.redisClient.Set(ctx, key, "1", s.refreshExpiry).Err()
	if err != nil {
		s.logger.Errorf("Failed to revoke user tokens: %v", err)
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	s.logger.Infof("All tokens revoked for user %s", userID)
	return nil
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *JWTService) IsTokenBlacklisted(ctx context.Context, tokenString string, tokenType string) (bool, error) {
	if s.redisClient == nil {
		return false, nil
	}

	tokenHash := s.hashToken(tokenString)
	key := fmt.Sprintf("blacklist:%s:%s", tokenType, tokenHash)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	return exists > 0, nil
}

// IsUserRevoked checks if a user's tokens are revoked
func (s *JWTService) IsUserRevoked(ctx context.Context, userID string) (bool, error) {
	if s.redisClient == nil {
		return false, nil
	}

	key := fmt.Sprintf("revoked_users:%s", userID)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check user revocation: %w", err)
	}

	return exists > 0, nil
}

// hashToken creates a SHA256 hash of the token for blacklisting
func (s *JWTService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
