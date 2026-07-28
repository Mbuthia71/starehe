package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/services"
	"starehian-society-platform/internal/tasks"
	"starehian-society-platform/pkg/config"
	"starehian-society-platform/pkg/database"
	"starehian-society-platform/pkg/logger"
	"starehian-society-platform/pkg/redis"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting Starehian Society Platform Worker...")

	// Initialize database
	db, err := database.NewDB(cfg.Database.URL)
	if err != nil {
		appLogger.Errorf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.NewRedis(cfg.Redis.URL)
	if err != nil {
		appLogger.Errorf("Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Initialize repositories
	pointsRepo := repository.NewPointsRepository(db)
	userRepo := repository.NewUserRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Initialize services
	pointsService := services.NewPointsService(pointsRepo, userRepo, db.DB, appLogger)

	// Initialize task handlers
	taskHandler := tasks.NewTaskHandler(pointsService, notificationRepo, appLogger)

	// Configure Asynq server
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.Redis.URL,
		// Add password if configured
		// Password: os.Getenv("REDIS_PASSWORD"),
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // Number of concurrent workers
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":     1,
			},
			RetryDelayFunc: func(n int, err error, task *asynq.Task) time.Duration {
				// Exponential backoff: 5s, 10s, 20s, 40s, 80s, max 5min
				delay := time.Duration(5*(1<<(n-1))) * time.Second
				if delay > 5*time.Minute {
					delay = 5 * time.Minute
				}
				return delay
			},
		},
	)

	// Register task handlers
	mux := asynq.NewServeMux()

	// Points redemption tasks
	mux.HandleFunc(tasks.TypeProcessRedemption, taskHandler.ProcessRedemption)
	mux.HandleFunc(tasks.TypeProcessPartnerCallback, taskHandler.ProcessPartnerCallback)

	// Email tasks
	mux.HandleFunc(tasks.TypeSendEmail, taskHandler.SendEmail)

	// Notification tasks
	mux.HandleFunc(tasks.TypeSendPushNotification, taskHandler.SendPushNotification)

	// Media processing tasks
	mux.HandleFunc(tasks.TypeProcessMedia, taskHandler.ProcessMedia)

	// Analytics tasks
	mux.HandleFunc(tasks.TypeAggregateAnalytics, taskHandler.AggregateAnalytics)

	// Start worker
	appLogger.Info("Worker started successfully")

	// Run server in a goroutine
	go func() {
		if err := srv.Run(mux); err != nil {
			appLogger.Errorf("Worker server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down worker...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Errorf("Worker shutdown error: %v", err)
	}

	appLogger.Info("Worker stopped")
}
