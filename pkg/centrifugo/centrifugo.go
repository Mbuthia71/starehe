package centrifugo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"starehian-society-platform/pkg/logger"
)

type CentrifugoClient struct {
	apiKey     string
	apiSecret  string
	hmacSecret string
	httpClient *http.Client
	baseURL    string
	logger     *logger.Logger
}

type PublishRequest struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type BroadcastRequest struct {
	Channels []string  `json:"channels"`
	Data     interface{} `json:"data"`
}

type PresenceRequest struct {
	Channel string `json:"channel"`
}

type PresenceResponse struct {
	Presence map[string]PresenceInfo `json:"presence"`
}

type PresenceInfo struct {
	ClientID string                 `json:"client"`
	UserID   string                 `json:"user"`
	UserInfo map[string]interface{} `json:"info"`
}

func NewCentrifugoClient(apiKey, apiSecret, hmacSecret, baseURL string, logger *logger.Logger) *CentrifugoClient {
	return &CentrifugoClient{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		hmacSecret: hmacSecret,
		baseURL:    baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GenerateConnectionToken generates a JWT token for Centrifugo connection
func (c *CentrifugoClient) GenerateConnectionToken(userID string, expire int64) (string, error) {
	now := time.Now().Unix()
	if expire == 0 {
		expire = now + 24*60*60 // 24 hours default
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expire,
		"iat": now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.hmacSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateChannelToken generates a JWT token for private channel subscription
func (c *CentrifugoClient) GenerateChannelToken(userID, channel string, expire int64) (string, error) {
	now := time.Now().Unix()
	if expire == 0 {
		expire = now + 24*60*60 // 24 hours default
	}

	claims := jwt.MapClaims{
		"sub":     userID,
		"channel": channel,
		"exp":     expire,
		"iat":     now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.hmacSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Publish publishes data to a specific channel
func (c *CentrifugoClient) Publish(ctx context.Context, channel string, data interface{}) error {
	req := PublishRequest{
		Channel: channel,
		Data:    data,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/publish", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("centrifugo returned status %d", resp.StatusCode)
	}

	c.logger.Infof("Published to channel %s", channel)
	return nil
}

// Broadcast publishes data to multiple channels
func (c *CentrifugoClient) Broadcast(ctx context.Context, channels []string, data interface{}) error {
	req := BroadcastRequest{
		Channels: channels,
		Data:     data,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/broadcast", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("centrifugo returned status %d", resp.StatusCode)
	}

	c.logger.Infof("Broadcasted to %d channels", len(channels))
	return nil
}

// Presence gets presence information for a channel
func (c *CentrifugoClient) Presence(ctx context.Context, channel string) (*PresenceResponse, error) {
	req := PresenceRequest{
		Channel: channel,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/presence", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("centrifugo returned status %d", resp.StatusCode)
	}

	var presenceResp PresenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&presenceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &presenceResp, nil
}

// GenerateHMACSignature generates HMAC signature for Centrifugo API
func (c *CentrifugoClient) GenerateHMACSignature(data string) string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
