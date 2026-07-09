package middleware

import (
	"context"
	"net/http"
	"strings"

	"starehian-society-platform/internal/auth"
	"starehian-society-platform/pkg/logger"
)

type AuthMiddleware struct {
	jwtService *auth.JWTService
	logger     *logger.Logger
}

func NewAuthMiddleware(jwtService *auth.JWTService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// Authenticate validates the JWT token and adds user context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			m.logger.Errorf("Failed to validate token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if the user has the required role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("user_role")
			if userRole == nil {
				http.Error(w, "User role not found in context", http.StatusUnauthorized)
				return
			}

			if userRole.(string) != role {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin checks if the user has any admin role
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("user_role")
		if userRole == nil {
			http.Error(w, "User role not found in context", http.StatusUnauthorized)
			return
		}

		role := userRole.(string)
		if role != "super_admin" && role != "moderator" && role != "support" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) string {
	if userRole, ok := ctx.Value("user_role").(string); ok {
		return userRole
	}
	return ""
}
