package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	AfricasTalking AfricasTalkingConfig
	R2       R2Config
	FCM      FCMConfig
	Centrifugo CentrifugoConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret           string
	RefreshSecret    string
	AccessTokenExpiry  int // in hours
	RefreshTokenExpiry int // in days
}

type AfricasTalkingConfig struct {
	APIKey   string
	Username string
}

type R2Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
}

type FCMConfig struct {
	ServerKey string
}

type CentrifugoConfig struct {
	Secret    string
	APIKey    string
	AdminPassword string
	AdminSecret string
	APIURL    string
}

type ServerConfig struct {
	Port string
	Environment string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://starehian_user:changeme@localhost:5432/starehian_db?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "changeme"),
			RefreshSecret:    getEnv("REFRESH_TOKEN_SECRET", "changeme"),
			AccessTokenExpiry: getEnvAsInt("JWT_ACCESS_EXPIRY_HOURS", 24),
			RefreshTokenExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		},
		AfricasTalking: AfricasTalkingConfig{
			APIKey:   getEnv("AFRICAS_TALKING_API_KEY", ""),
			Username: getEnv("AFRICAS_TALKING_USERNAME", ""),
		},
		R2: R2Config{
			AccessKey: getEnv("CLOUDFLARE_R2_ACCESS_KEY", ""),
			SecretKey: getEnv("CLOUDFLARE_R2_SECRET_KEY", ""),
			Bucket:    getEnv("CLOUDFLARE_R2_BUCKET", ""),
			Endpoint:  getEnv("CLOUDFLARE_R2_ENDPOINT", ""),
		},
		FCM: FCMConfig{
			ServerKey: getEnv("FCM_SERVER_KEY", ""),
		},
		Centrifugo: CentrifugoConfig{
			Secret:    getEnv("CENTRIFUGO_SECRET", "changeme"),
			APIKey:    getEnv("CENTRIFUGO_API_KEY", "changeme"),
			AdminPassword: getEnv("CENTRIFUGO_ADMIN_PASSWORD", "changeme"),
			AdminSecret: getEnv("CENTRIFUGO_ADMIN_SECRET", "changeme"),
			APIURL:    getEnv("CENTRIFUGO_API_URL", "http://localhost:8000"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
