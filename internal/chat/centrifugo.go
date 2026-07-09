package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"starehian-society-platform/pkg/config"
)

type CentrifugoClient struct {
	apiKey    string
	apiSecret string
	apiURL    string
	client    *http.Client
}

type CentrifugoMessage struct {
	Data interface{} `json:"data"`
}

type CentrifugoBroadcastRequest struct {
	Channels []string             `json:"channels"`
	Data     interface{}          `json:"data"`
}

type CentrifugoPublishRequest struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type CentrifugoPresenceRequest struct {
	Channel string `json:"channel"`
}

type CentrifugoPresenceResponse struct {
	Presence map[string]CentrifugoClientInfo `json:"presence"`
}

type CentrifugoClientInfo struct {
	User   string `json:"user"`
	Client string `json:"client"`
}

func NewCentrifugoClient(cfg *config.CentrifugoConfig) *CentrifugoClient {
	return &CentrifugoClient{
		apiKey:    cfg.APIKey,
		apiSecret: cfg.Secret,
		apiURL:    cfg.APIURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Broadcast sends a message to multiple channels
func (c *CentrifugoClient) Broadcast(ctx context.Context, channels []string, data interface{}) error {
	url := fmt.Sprintf("%s/api/broadcast", c.apiURL)
	
	req := CentrifugoBroadcastRequest{
		Channels: channels,
		Data:     data,
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)
	
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send broadcast: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("broadcast failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// Publish sends a message to a single channel
func (c *CentrifugoClient) Publish(ctx context.Context, channel string, data interface{}) error {
	url := fmt.Sprintf("%s/api/publish", c.apiURL)
	
	req := CentrifugoPublishRequest{
		Channel: channel,
		Data:    data,
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)
	
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("publish failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// GetPresence gets online users in a channel
func (c *CentrifugoClient) GetPresence(ctx context.Context, channel string) (*CentrifugoPresenceResponse, error) {
	url := fmt.Sprintf("%s/api/presence", c.apiURL)
	
	req := CentrifugoPresenceRequest{
		Channel: channel,
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "apikey "+c.apiKey)
	
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("presence failed with status %d", resp.StatusCode)
	}
	
	var presenceResp CentrifugoPresenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&presenceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &presenceResp, nil
}

// GenerateConnectionToken generates a token for a user to connect to Centrifugo
func (c *CentrifugoClient) GenerateConnectionToken(userID string) (string, error) {
	// In production, use proper JWT signing with the Centrifugo secret
	// For now, return a placeholder
	return fmt.Sprintf("token_%s", userID), nil
}

// GenerateChannelToken generates a token for a user to subscribe to a private channel
func (c *CentrifugoClient) GenerateChannelToken(userID, channel string) (string, error) {
	// In production, use proper JWT signing
	return fmt.Sprintf("channel_token_%s_%s", userID, channel), nil
}
