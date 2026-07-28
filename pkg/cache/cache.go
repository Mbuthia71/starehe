package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"starehian-society-platform/pkg/logger"
)

type CacheService struct {
	client *redis.Client
	logger *logger.Logger
}

func NewCacheService(client *redis.Client, appLogger *logger.Logger) *CacheService {
	return &CacheService{
		client: client,
		logger: appLogger,
	}
}

// CacheResult represents a cached result with metadata
type CacheResult struct {
	Data      interface{} `json:"data"`
	Version   string      `json:"version"`
	CachedAt  time.Time   `json:"cached_at"`
	TTL       time.Duration `json:"ttl"`
}

// Get retrieves a value from cache
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Cache miss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return nil
}

// Set stores a value in cache with TTL
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// DeletePattern removes all keys matching a pattern
func (s *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			s.logger.Errorf("Failed to delete key %s: %v", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}
	return nil
}

// InvalidateUser invalidates all cache entries for a specific user
func (s *CacheService) InvalidateUser(ctx context.Context, userID string) error {
	patterns := []string{
		fmt.Sprintf("profile:%s:*", userID),
		fmt.Sprintf("feed:%s:*", userID),
		fmt.Sprintf("posts:user:%s:*", userID),
		fmt.Sprintf("connections:%s:*", userID),
	}

	for _, pattern := range patterns {
		if err := s.DeletePattern(ctx, pattern); err != nil {
			s.logger.Errorf("Failed to invalidate pattern %s: %v", pattern, err)
		}
	}

	s.logger.Infof("Invalidated cache for user %s", userID)
	return nil
}

// InvalidatePost invalidates cache entries related to a post
func (s *CacheService) InvalidatePost(ctx context.Context, postID string) error {
	patterns := []string{
		fmt.Sprintf("post:%s", postID),
		fmt.Sprintf("feed:*"), // Invalidate all feeds as they may contain this post
	}

	for _, pattern := range patterns {
		if err := s.DeletePattern(ctx, pattern); err != nil {
			s.logger.Errorf("Failed to invalidate pattern %s: %v", pattern, err)
		}
	}

	s.logger.Infof("Invalidated cache for post %s", postID)
	return nil
}

// GenerateKey generates a cache key from components
func (s *CacheService) GenerateKey(prefix string, components ...string) string {
	key := prefix
	for _, comp := range components {
		if comp != "" {
			key += ":" + comp
		}
	}
	return key
}

// GenerateHashKey generates a hash-based cache key for complex queries
func (s *CacheService) GenerateHashKey(prefix string, params map[string]interface{}) string {
	hash := sha256.New()
	for k, v := range params {
		hash.Write([]byte(k))
		hash.Write([]byte(fmt.Sprintf("%v", v)))
	}
	hashStr := hex.EncodeToString(hash.Sum(nil))[:16]
	return fmt.Sprintf("%s:%s", prefix, hashStr)
}

// GetOrSet retrieves from cache or executes the function to get and cache the result
func (s *CacheService) GetOrSet(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	// Try to get from cache first
	if err := s.Get(ctx, key, dest); err == nil {
		return nil // Cache hit
	}

	// Cache miss, execute function
	result, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := s.Set(ctx, key, result, ttl); err != nil {
		s.logger.Warnf("Failed to cache result: %v", err)
		// Don't return error, the function succeeded
	}

	// Set dest to result
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return json.Unmarshal(resultBytes, dest)
}

// Cache profile with versioning
func (s *CacheService) CacheProfile(ctx context.Context, userID string, profile interface{}, ttl time.Duration) error {
	key := s.GenerateKey("profile", userID)
	return s.Set(ctx, key, profile, ttl)
}

// GetProfile retrieves cached profile
func (s *CacheService) GetProfile(ctx context.Context, userID string, dest interface{}) error {
	key := s.GenerateKey("profile", userID)
	return s.Get(ctx, key, dest)
}

// Cache feed with versioning
func (s *CacheService) CacheFeed(ctx context.Context, userID string, params map[string]interface{}, feed interface{}, ttl time.Duration) error {
	key := s.GenerateHashKey("feed", map[string]interface{}{"user_id": userID})
	for k, v := range params {
		key = s.GenerateHashKey(key, map[string]interface{}{k: v})
	}
	return s.Set(ctx, key, feed, ttl)
}

// GetFeed retrieves cached feed
func (s *CacheService) GetFeed(ctx context.Context, userID string, params map[string]interface{}, dest interface{}) error {
	key := s.GenerateHashKey("feed", map[string]interface{}{"user_id": userID})
	for k, v := range params {
		key = s.GenerateHashKey(key, map[string]interface{}{k: v})
	}
	return s.Get(ctx, key, dest)
}

// Cache search results
func (s *CacheService) CacheSearch(ctx context.Context, params map[string]interface{}, results interface{}, ttl time.Duration) error {
	key := s.GenerateHashKey("search", params)
	return s.Set(ctx, key, results, ttl)
}

// GetSearch retrieves cached search results
func (s *CacheService) GetSearch(ctx context.Context, params map[string]interface{}, dest interface{}) error {
	key := s.GenerateHashKey("search", params)
	return s.Get(ctx, key, dest)
}

// GetStats returns cache statistics
func (s *CacheService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := s.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Parse key count
	dbSize, err := s.client.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get db size: %w", err)
	}

	return map[string]interface{}{
		"info":    info,
		"db_size": dbSize,
	}, nil
}

// FlushAll clears all cache entries (use with caution)
func (s *CacheService) FlushAll(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	s.logger.Warn("Cache flushed")
	return nil
}
