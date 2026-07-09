package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"starehian-society-platform/pkg/config"
)

type AfricasTalkingService struct {
	apiKey   string
	username string
	client   *http.Client
}

type AfricasTalkingResponse struct {
	SMSMessageData struct {
		Recipients []struct {
			MessageID      string `json:"messageId"`
			Number         string `json:"number"`
			Status         string `json:"status"`
			Cost           string `json:"cost"`
			MessagePartID  int    `json:"messagePartId"`
		} `json:"Recipients"`
	} `json:"SMSMessageData"`
}

type AfricasTalkingError struct {
	Message string `json:"message"`
}

func NewAfricasTalkingService(cfg *config.AfricasTalkingConfig) *AfricasTalkingService {
	return &AfricasTalkingService{
		apiKey:   cfg.APIKey,
		username: cfg.Username,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendOTP sends an OTP code to the given phone number
func (s *AfricasTalkingService) SendOTP(phone, otp string) error {
	message := fmt.Sprintf("Your Starehian Society verification code is: %s. This code expires in 10 minutes.", otp)
	
	url := "https://api.africastalking.com/version1/messaging"
	
	data := map[string]string{
		"username": s.username,
		"to":       phone,
		"message":  message,
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", s.apiKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		var errResp AfricasTalkingError
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("OTP send failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("OTP send failed: %s", errResp.Message)
	}
	
	var respData AfricasTalkingResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	if len(respData.SMSMessageData.Recipients) == 0 {
		return fmt.Errorf("no recipients in response")
	}
	
	// Check if message was sent successfully
	for _, recipient := range respData.SMSMessageData.Recipients {
		if recipient.Status != "Success" {
			return fmt.Errorf("OTP send failed for %s: %s", recipient.Number, recipient.Status)
		}
	}
	
	return nil
}

// FormatPhone formats a phone number to international format
func FormatPhone(phone string) string {
	// Remove all non-digit characters
	formatted := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			formatted += string(char)
		}
	}
	
	// Add country code if not present (assuming Kenya: +254)
	if len(formatted) == 10 && formatted[0] == '0' {
		formatted = "254" + formatted[1:]
	} else if len(formatted) == 9 {
		formatted = "254" + formatted
	}
	
	return "+" + formatted
}
