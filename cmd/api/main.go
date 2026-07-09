package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"starehian-society-platform/internal/admin"
	"starehian-society-platform/internal/auth"
	"starehian-society-platform/internal/chat"
	"starehian-society-platform/internal/connections"
	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/notifications"
	"starehian-society-platform/internal/posts"
	"starehian-society-platform/internal/profiles"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/services"
	"starehian-society-platform/internal/tokens"
	"starehian-society-platform/pkg/config"
	"starehian-society-platform/pkg/database"
	"starehian-society-platform/pkg/logger"
	"starehian-society-platform/pkg/ratelimit"
	"starehian-society-platform/pkg/redis"
	"starehian-society-platform/pkg/storage"
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
	appLogger.Info("Starting Starehian Society Platform API...")

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

	// Initialize rate limiter
	rateLimiter := ratelimit.NewRateLimiter(redisClient)

	// Initialize JWT service
	jwtService := tokens.NewJWTService(&cfg.JWT)

	// Initialize Africa's Talking service
	atService := tokens.NewAfricasTalkingService(&cfg.AfricasTalking)

	// Initialize OTP service
	otpService := tokens.NewOTPService(redisClient, rateLimiter, atService, appLogger)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	connectionRepo := repository.NewConnectionRepository(db)
	postRepo := repository.NewPostRepository(db)
	chatRepo := repository.NewChatRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Initialize authorization service
	authzService := middleware.NewAuthorizationService(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, profileRepo, redisClient, jwtService, otpService, appLogger)
	profileService := services.NewProfileService(profileRepo, userRepo, appLogger)
	adminService := services.NewAdminService(userRepo, adminRepo, appLogger)
	connectionService := services.NewConnectionService(connectionRepo, userRepo, appLogger)
	postService := services.NewPostService(postRepo, appLogger)
	chatService := services.NewChatService(chatRepo, connectionRepo, redisClient, appLogger)
	analyticsService := services.NewAnalyticsService(db, appLogger)
	broadcastService := services.NewBroadcastService(userRepo, profileRepo, notificationRepo, appLogger)

	// Initialize storage
	r2Storage := storage.NewR2Storage(&cfg.R2)

	// Initialize handlers
	authHandler := auth.NewAuthHandler(authService, rateLimiter)
	profileHandler := profiles.NewProfileHandler(profileService)
	adminHandler := admin.NewAdminHandler(adminService, userRepo, authzService)
	connectionHandler := connections.NewConnectionHandler(connectionService)
	postHandler := posts.NewPostHandler(postService)
	uploadHandler := posts.NewUploadHandler(r2Storage)
	chatHandler := chat.NewChatHandler(chatService)
	notificationHandler := notifications.NewNotificationHandler(notificationRepo)
	analyticsHandler := admin.NewAnalyticsHandler(analyticsService)
	broadcastHandler := admin.NewBroadcastHandler(broadcastService)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, appLogger)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database
		dbStatus := "healthy"
		if err := db.Ping(); err != nil {
			dbStatus = "unhealthy"
		}

		// Check Redis
		redisStatus := "healthy"
		if err := redisClient.Ping(c.Context()); err != nil {
			redisStatus = "unhealthy"
		}

		overallStatus := "healthy"
		if dbStatus != "healthy" || redisStatus != "healthy" {
			overallStatus = "degraded"
		}

		return c.JSON(fiber.Map{
			"status":   overallStatus,
			"database": dbStatus,
			"redis":    redisStatus,
		})
	})

	// Readiness check
	app.Get("/ready", func(c *fiber.Ctx) error {
		// Check if all dependencies are ready
		if err := db.Ping(); err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "database not ready",
			})
		}

		if err := redisClient.Ping(c.Context()); err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "redis not ready",
			})
		}

		return c.JSON(fiber.Map{
			"status": "ready",
		})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes (public) with rate limiting
	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.RateLimitMiddleware(rateLimiter, middleware.AuthRateLimit))
	authRoutes.Post("/request-otp", authHandler.RequestOTP)
	authRoutes.Post("/signup", authHandler.Signup)
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Post("/admin/login", authHandler.AdminLogin)
	authRoutes.Post("/refresh", authHandler.RefreshToken)
	authRoutes.Post("/logout", authHandler.Logout)

	// Protected routes with general rate limiting
	protected := api.Group("/")
	protected.Use(authMiddleware.Authenticate)
	protected.Use(middleware.RateLimitMiddleware(rateLimiter, middleware.GeneralRateLimit))

	// Profile routes
	profileRoutes := protected.Group("/profiles")
	profileRoutes.Get("/me", profileHandler.GetOwnProfile)
	profileRoutes.Put("/me", profileHandler.UpdateProfile)
	profileRoutes.Get("/:id", profileHandler.GetProfile)
	profileRoutes.Post("/search", profileHandler.SearchProfiles)

	// Connection routes
	connectionRoutes := protected.Group("/connections")
	connectionRoutes.Post("/", connectionHandler.SendConnectionRequest)
	connectionRoutes.Post("/:id/accept", connectionHandler.AcceptConnectionRequest)
	connectionRoutes.Post("/:id/reject", connectionHandler.RejectConnectionRequest)
	connectionRoutes.Delete("/:id", connectionHandler.RemoveConnection)
	connectionRoutes.Post("/:id/block", connectionHandler.BlockUser)
	connectionRoutes.Delete("/:id/block", connectionHandler.UnblockUser)
	connectionRoutes.Get("/", connectionHandler.GetConnections)
	connectionRoutes.Get("/pending", connectionHandler.GetPendingConnections)
	connectionRoutes.Get("/sent", connectionHandler.GetSentConnections)
	connectionRoutes.Get("/blocked", connectionHandler.GetBlocks)

	// Post routes
	postRoutes := protected.Group("/posts")
	postRoutes.Post("/", postHandler.CreatePost)
	postRoutes.Get("/feed", postHandler.GetFeed)
	postRoutes.Get("/:id", postHandler.GetPost)
	postRoutes.Put("/:id", postHandler.UpdatePost)
	postRoutes.Delete("/:id", postHandler.DeletePost)
	postRoutes.Get("/user/:id", postHandler.GetPostsByUser)
	postRoutes.Post("/:id/comments", postHandler.CreateComment)
	postRoutes.Get("/:id/comments", postHandler.GetComments)
	postRoutes.Delete("/:id/comments/:commentId", postHandler.DeleteComment)
	postRoutes.Post("/:id/reactions", postHandler.CreateReaction)
	postRoutes.Get("/:id/reactions", postHandler.GetReactions)
	postRoutes.Delete("/:id/reactions", postHandler.DeleteReaction)
	postRoutes.Get("/:id/reaction", postHandler.GetUserReaction)

	// Media upload routes with stricter rate limiting
	uploadRoutes := protected.Group("/upload")
	uploadRoutes.Use(middleware.RateLimitMiddleware(rateLimiter, middleware.UploadRateLimit))
	uploadRoutes.Post("/media", uploadHandler.UploadMedia)
	uploadRoutes.Post("/media/multiple", uploadHandler.UploadMultipleMedia)

	// Chat routes
	chatRoutes := protected.Group("/chat")
	chatRoutes.Post("/direct/:id", chatHandler.CreateDirectConversation)
	chatRoutes.Post("/group", chatHandler.CreateGroupConversation)
	chatRoutes.Post("/:id/messages", chatHandler.SendMessage)
	chatRoutes.Get("/:id/messages", chatHandler.GetMessages)
	chatRoutes.Get("/", chatHandler.GetConversations)
	chatRoutes.Post("/:id/read", chatHandler.MarkAsRead)
	chatRoutes.Post("/:id/members", chatHandler.AddMember)
	chatRoutes.Delete("/:id/members/:memberId", chatHandler.RemoveMember)
	chatRoutes.Delete("/:id", chatHandler.LeaveConversation)
	chatRoutes.Get("/token", chatHandler.GetConnectionToken)
	chatRoutes.Get("/channel-token", chatHandler.GetChannelToken)

	// Notification routes
	notificationRoutes := protected.Group("/notifications")
	notificationRoutes.Get("/", notificationHandler.GetNotifications)
	notificationRoutes.Get("/unread-count", notificationHandler.GetUnreadCount)
	notificationRoutes.Post("/:id/read", notificationHandler.MarkAsRead)
	notificationRoutes.Post("/read-all", notificationHandler.MarkAllAsRead)
	notificationRoutes.Delete("/:id", notificationHandler.DeleteNotification)

	// Admin routes (admin only) with admin rate limiting
	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(authMiddleware.RequireAdmin)
	adminRoutes.Use(middleware.RateLimitMiddleware(rateLimiter, middleware.AdminRateLimit))

	// User management
	adminRoutes.Get("/users", adminHandler.GetUsers)
	adminRoutes.Get("/users/:id", adminHandler.GetUser)
	adminRoutes.Post("/users/:id/suspend", adminHandler.SuspendUser)
	adminRoutes.Post("/users/:id/activate", adminHandler.ActivateUser)
	adminRoutes.Put("/users/:id/role", adminHandler.UpdateUserRole)

	// Reports
	adminRoutes.Post("/reports", adminHandler.CreateReport)
	adminRoutes.Get("/reports", adminHandler.GetReports)
	adminRoutes.Put("/reports/:id", adminHandler.UpdateReportStatus)

	// Audit log
	adminRoutes.Get("/audit", adminHandler.GetAdminActions)

	// Alumni roster
	adminRoutes.Post("/roster", adminHandler.AddAlumniRosterEntry)
	adminRoutes.Get("/roster", adminHandler.GetAlumniRoster)
	adminRoutes.Post("/users/:id/verify", adminHandler.VerifyUserAgainstRoster)

	// Analytics
	adminRoutes.Get("/analytics/dashboard", analyticsHandler.GetDashboardMetrics)
	adminRoutes.Get("/analytics/cohorts", analyticsHandler.GetCohortAnalytics)
	adminRoutes.Get("/analytics/signups", analyticsHandler.GetSignupsOverTime)
	adminRoutes.Get("/analytics/engagement", analyticsHandler.GetEngagementOverTime)
	adminRoutes.Get("/analytics/top-content", analyticsHandler.GetTopContent)

	// Broadcasts
	adminRoutes.Post("/broadcasts", broadcastHandler.SendBroadcast)

	// TODO: Re-enable bulk operations when BulkOperationService is implemented
	// // Bulk operations
	// adminRoutes.Post("/bulk/suspend", bulkHandler.BulkSuspend)
	// adminRoutes.Post("/bulk/activate", bulkHandler.BulkActivate)
	// adminRoutes.Post("/bulk/verify", bulkHandler.BulkVerify)
	// adminRoutes.Post("/bulk/delete", bulkHandler.BulkDelete)

	// Start server
	port := cfg.Server.Port
	appLogger.Infof("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		appLogger.Errorf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
