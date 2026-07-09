package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"starehian-society-platform/pkg/config"
)

type R2Storage struct {
	accessKey string
	secretKey string
	bucket    string
	endpoint  string
	client    *http.Client
}

func NewR2Storage(cfg *config.R2Config) *R2Storage {
	return &R2Storage{
		accessKey: cfg.AccessKey,
		secretKey: cfg.SecretKey,
		bucket:    cfg.Bucket,
		endpoint:  cfg.Endpoint,
		client:    &http.Client{},
	}
}

// UploadFile uploads a file to R2
func (s *R2Storage) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	// For simplicity, we'll use a basic HTTP upload
	// In production, you'd use the AWS S3 SDK with R2 endpoint
	
	url := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key)
	
	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "PUT", url, reader)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Amz-Content-Sha256", "UNSIGNED-PAYLOAD")
	
	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("upload failed with status %d", resp.StatusCode)
	}
	
	// Return the public URL
	publicURL := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(s.endpoint, "/"), s.bucket, key)
	return publicURL, nil
}

// DeleteFile deletes a file from R2
func (s *R2Storage) DeleteFile(ctx context.Context, key string) error {
	url := fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// GenerateKey generates a unique key for a file
func (s *R2Storage) GenerateKey(userID, fileType string) string {
	// Simple key generation: userID/timestamp/type
	// In production, use UUID or timestamp
	return fmt.Sprintf("%s/%s", userID, fileType)
}
