package centrifugo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CentrifugoClient wraps Centrifugo API interactions
type CentrifugoClient struct {
	apiKey string
	secret string
}

// NewCentrifugoClient creates a new Centrifugo client
func NewCentrifugoClient(apiKey, secret string) (*CentrifugoClient, error) {
	if apiKey == "" || secret == "" {
		return nil, fmt.Errorf("centrifugo API key and secret are required")
	}

	return &CentrifugoClient{
		apiKey: apiKey,
		secret: secret,
	}, nil
}

// GenerateConnectionToken generates a token for WebSocket connection
func (c *CentrifugoClient) GenerateConnectionToken(userID string, expiration int64) (string, error) {
	if expiration == 0 {
		expiration = 7200 // 2 hours default
	}

	now := time.Now().Unix()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now,
		"exp": now + expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateChannelToken generates a token for channel subscription
func (c *CentrifugoClient) GenerateChannelToken(
	userID string,
	channel string,
	expiration int64,
) (string, error) {
	if expiration == 0 {
		expiration = 7200 // 2 hours default
	}

	now := time.Now().Unix()

	claims := jwt.MapClaims{
		"sub":     userID,
		"channel": channel,
		"iat":     now,
		"exp":     now + expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// Publish publishes data to a channel (for server-side broadcasting)
func (c *CentrifugoClient) Publish(
	ctx context.Context,
	channel string,
	data interface{},
) error {
	// Note: This would typically require an HTTP client to call Centrifugo's HTTP API.
	// For now, this is a placeholder that demonstrates the structure.
	// In production, you would implement the actual HTTP call.

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// TODO: Implement actual HTTP call to Centrifugo API endpoint
	// POST to http://centrifugo:8000/api/publish
	// with Authorization header and request body containing channel and data

	_ = dataJSON // Use the variable to avoid linting error
	return nil
}
