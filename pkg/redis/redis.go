package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func NewRedis(redisURL string) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &Redis{client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

// Session management
func (r *Redis) SetSession(ctx context.Context, sessionID string, userID string, expiry time.Duration) error {
	return r.Set(ctx, "session:"+sessionID, userID, expiry).Err()
}

func (r *Redis) GetSession(ctx context.Context, sessionID string) (string, error) {
	return r.Get(ctx, "session:"+sessionID).Result()
}

func (r *Redis) DeleteSession(ctx context.Context, sessionID string) error {
	return r.Del(ctx, "session:"+sessionID).Err()
}

// OTP storage
func (r *Redis) SetOTP(ctx context.Context, phone string, otp string, expiry time.Duration) error {
	return r.Set(ctx, "otp:"+phone, otp, expiry).Err()
}

func (r *Redis) GetOTP(ctx context.Context, phone string) (string, error) {
	return r.Get(ctx, "otp:"+phone).Result()
}

func (r *Redis) DeleteOTP(ctx context.Context, phone string) error {
	return r.Del(ctx, "otp:"+phone).Err()
}

// Rate limiting
func (r *Redis) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := r.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return incr.Val(), nil
}

func (r *Redis) GetRateLimitCount(ctx context.Context, key string) (int64, error) {
	return r.Get(ctx, key).Int64()
}

// Online presence
func (r *Redis) SetUserOnline(ctx context.Context, userID string, expiry time.Duration) error {
	return r.Set(ctx, "online:"+userID, "1", expiry).Err()
}

func (r *Redis) SetUserOffline(ctx context.Context, userID string) error {
	return r.Del(ctx, "online:"+userID).Err()
}

func (r *Redis) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	result, err := r.Exists(ctx, "online:"+userID).Result()
	return result > 0, err
}
