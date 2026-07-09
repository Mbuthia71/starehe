package middleware

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type AuthorizationService struct {
	db *database.DB
}

func NewAuthorizationService(db *database.DB) *AuthorizationService {
	return &AuthorizationService{db: db}
}

// CanViewProfile checks if a viewer can view a target's profile
func (a *AuthorizationService) CanViewProfile(ctx context.Context, viewerID, targetID string) (bool, error) {
	if viewerID == targetID {
		return true, nil // Users can always view their own profile
	}

	// Check if viewer has blocked target or vice versa
	blocked, err := a.checkBlock(ctx, viewerID, targetID)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	// Get target's profile visibility
	var visibility string
	query := `SELECT profile_visibility FROM profiles WHERE user_id = $1`
	err = a.db.GetContext(ctx, &visibility, query, targetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Profile doesn't exist
		}
		return false, fmt.Errorf("failed to get profile visibility: %w", err)
	}

	switch visibility {
	case string(models.VisibilityPublic):
		return true, nil
	case string(models.VisibilityPrivate):
		return false, nil
	case string(models.VisibilityConnections):
		// Check if they are connected
		connected, err := a.checkConnection(ctx, viewerID, targetID)
		if err != nil {
			return false, err
		}
		return connected, nil
	default:
		return false, nil
	}
}

// CanViewProfileSection checks if a viewer can view a specific section of a profile
func (a *AuthorizationService) CanViewProfileSection(ctx context.Context, viewerID, targetID, section string) (bool, error) {
	if viewerID == targetID {
		return true, nil
	}

	// Check if viewer has blocked target or vice versa
	blocked, err := a.checkBlock(ctx, viewerID, targetID)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	// Get the section's visibility
	var visibility string
	var query string
	switch section {
	case "contact":
		query = `SELECT contact_visibility FROM profiles WHERE user_id = $1`
	case "career":
		query = `SELECT career_visibility FROM profiles WHERE user_id = $1`
	default:
		// For other sections, use profile visibility
		query = `SELECT profile_visibility FROM profiles WHERE user_id = $1`
	}

	err = a.db.GetContext(ctx, &visibility, query, targetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to get section visibility: %w", err)
	}

	switch visibility {
	case string(models.VisibilityPublic):
		return true, nil
	case string(models.VisibilityPrivate):
		return false, nil
	case string(models.VisibilityConnections):
		connected, err := a.checkConnection(ctx, viewerID, targetID)
		if err != nil {
			return false, err
		}
		return connected, nil
	default:
		return false, nil
	}
}

// CanViewPost checks if a viewer can view a post
func (a *AuthorizationService) CanViewPost(ctx context.Context, viewerID, postID string) (bool, error) {
	// Get post visibility and author
	var visibility, authorID string
	query := `SELECT visibility, user_id FROM posts WHERE id = $1`
	err := a.db.GetContext(ctx, &visibility, &authorID, query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to get post: %w", err)
	}

	if viewerID == authorID {
		return true, nil // Author can always view their own posts
	}

	// Check if viewer has blocked author or vice versa
	blocked, err := a.checkBlock(ctx, viewerID, authorID)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	switch visibility {
	case string(models.VisibilityPublic):
		return true, nil
	case string(models.VisibilityPrivate):
		return false, nil
	case string(models.VisibilityConnections):
		connected, err := a.checkConnection(ctx, viewerID, authorID)
		if err != nil {
			return false, err
		}
		return connected, nil
	default:
		return false, nil
	}
}

// CanSendConnectionRequest checks if a viewer can send a connection request to target
func (a *AuthorizationService) CanSendConnectionRequest(ctx context.Context, viewerID, targetID string) (bool, error) {
	if viewerID == targetID {
		return false, nil // Cannot send request to self
	}

	// Check if blocked
	blocked, err := a.checkBlock(ctx, viewerID, targetID)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	// Check if already connected or pending
	var status string
	query := `SELECT status FROM connections WHERE user_id = $1 AND connected_user_id = $2`
	err = a.db.GetContext(ctx, &status, query, viewerID, targetID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check connection status: %w", err)
	}

	if status == string(models.ConnectionStatusAccepted) || status == string(models.ConnectionStatusPending) {
		return false, nil // Already connected or request pending
	}

	return true, nil
}

// checkBlock checks if either user has blocked the other
func (a *AuthorizationService) checkBlock(ctx context.Context, userID1, userID2 string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM blocks 
		WHERE (blocker_id = $1 AND blocked_id = $2) 
		OR (blocker_id = $2 AND blocked_id = $1)
	`
	err := a.db.GetContext(ctx, &count, query, userID1, userID2)
	if err != nil {
		return false, fmt.Errorf("failed to check block: %w", err)
	}
	return count > 0, nil
}

// checkConnection checks if two users are connected
func (a *AuthorizationService) checkConnection(ctx context.Context, userID1, userID2 string) (bool, error) {
	var status string
	query := `
		SELECT status FROM connections 
		WHERE (user_id = $1 AND connected_user_id = $2 AND status = 'accepted')
		OR (user_id = $2 AND connected_user_id = $1 AND status = 'accepted')
	`
	err := a.db.GetContext(ctx, &status, query, userID1, userID2)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check connection: %w", err)
	}
	return status == string(models.ConnectionStatusAccepted), nil
}

// HasRole checks if a user has a specific role
func (a *AuthorizationService) HasRole(ctx context.Context, userID string, role models.UserRole) (bool, error) {
	var userRole string
	query := `SELECT role FROM users WHERE id = $1`
	err := a.db.GetContext(ctx, &userRole, query, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user role: %w", err)
	}
	return userRole == string(role), nil
}

// IsAdmin checks if a user has any admin role
func (a *AuthorizationService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	var role string
	query := `SELECT role FROM users WHERE id = $1`
	err := a.db.GetContext(ctx, &role, query, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user role: %w", err)
	}
	return role == string(models.RoleSuperAdmin) || role == string(models.RoleModerator) || role == string(models.RoleSupport), nil
}

// IsSuperAdmin checks if a user is a super admin
func (a *AuthorizationService) IsSuperAdmin(ctx context.Context, userID string) (bool, error) {
	var role string
	query := `SELECT role FROM users WHERE id = $1`
	err := a.db.GetContext(ctx, &role, query, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user role: %w", err)
	}
	return role == string(models.RoleSuperAdmin), nil
}
