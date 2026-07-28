package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type ProfileRepository struct {
	db *database.DB
}

func NewProfileRepository(db *database.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create creates a new profile
func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	query := `
		INSERT INTO profiles (user_id, full_name, bio, avatar_url, cover_url, class_year, house, career, location, profile_visibility, contact_visibility, career_visibility)
		VALUES (:user_id, :full_name, :bio, :avatar_url, :cover_url, :class_year, :house, :career, :location, :profile_visibility, :contact_visibility, :career_visibility)
		RETURNING created_at, updated_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, profile)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&profile.CreatedAt, &profile.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan profile: %w", err)
		}
	}
	
	return nil
}

// GetByID retrieves a profile by user ID
func (r *ProfileRepository) GetByID(ctx context.Context, userID string) (*models.Profile, error) {
	var profile models.Profile
	query := `
		SELECT user_id, full_name, bio, avatar_url, cover_url, class_year, house, career, location, 
		       profile_visibility, contact_visibility, career_visibility, created_at, updated_at
		FROM profiles WHERE user_id = $1
	`
	
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	
	return &profile, nil
}

// Update updates a profile
func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	query := `
		UPDATE profiles 
		SET full_name = :full_name, bio = :bio, avatar_url = :avatar_url, cover_url = :cover_url, 
		    class_year = :class_year, house = :house, career = :career, location = :location,
		    profile_visibility = :profile_visibility, contact_visibility = :contact_visibility, career_visibility = :career_visibility
		WHERE user_id = :user_id
		RETURNING updated_at
	`
	
	result, err := r.db.NamedExecContext(ctx, query, profile)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("profile not found")
	}
	
	return nil
}

// UpdateVisibility updates profile visibility settings
func (r *ProfileRepository) UpdateVisibility(ctx context.Context, userID string, profileVisibility, contactVisibility, careerVisibility string) error {
	query := `
		UPDATE profiles 
		SET profile_visibility = $1, contact_visibility = $2, career_visibility = $3
		WHERE user_id = $4
	`
	
	result, err := r.db.ExecContext(ctx, query, profileVisibility, contactVisibility, careerVisibility, userID)
	if err != nil {
		return fmt.Errorf("failed to update profile visibility: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("profile not found")
	}
	
	return nil
}

// Search searches for profiles by various criteria using offset pagination
func (r *ProfileRepository) Search(ctx context.Context, searchTerm string, classYear *int, house, location, career string, limit, offset int) ([]*models.Profile, error) {
	var profiles []*models.Profile

	baseQuery := `
		SELECT p.user_id, p.full_name, p.bio, p.avatar_url, p.cover_url, p.class_year, p.house, p.career, p.location,
		       p.profile_visibility, p.contact_visibility, p.career_visibility, p.created_at, p.updated_at,
		       u.file_number
		FROM profiles p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if searchTerm != "" {
		baseQuery += fmt.Sprintf(" AND (p.full_name ILIKE $%d OR p.career ILIKE $%d OR p.location ILIKE $%d)", argCount, argCount+1, argCount+2)
		args = append(args, "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
		argCount += 3
	}

	if classYear != nil {
		baseQuery += fmt.Sprintf(" AND p.class_year = $%d", argCount)
		args = append(args, *classYear)
		argCount++
	}

	if house != "" {
		baseQuery += fmt.Sprintf(" AND p.house = $%d", argCount)
		args = append(args, house)
		argCount++
	}

	if location != "" {
		baseQuery += fmt.Sprintf(" AND p.location ILIKE $%d", argCount)
		args = append(args, "%"+location+"%")
		argCount++
	}

	if career != "" {
		baseQuery += fmt.Sprintf(" AND p.career ILIKE $%d", argCount)
		args = append(args, "%"+career+"%")
		argCount++
	}

	baseQuery += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	query, args, err := sqlx.In(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &profiles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}

	return profiles, nil
}

// SearchCursor searches for profiles by various criteria using cursor-based pagination
func (r *ProfileRepository) SearchCursor(ctx context.Context, searchTerm string, classYear *int, house, location, career string, limit int, cursor *string) ([]*models.Profile, *string, error) {
	var profiles []*models.Profile

	baseQuery := `
		SELECT p.user_id, p.full_name, p.bio, p.avatar_url, p.cover_url, p.class_year, p.house, p.career, p.location,
		       p.profile_visibility, p.contact_visibility, p.career_visibility, p.created_at, p.updated_at,
		       u.file_number
		FROM profiles p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if searchTerm != "" {
		baseQuery += fmt.Sprintf(" AND (p.full_name ILIKE $%d OR p.career ILIKE $%d OR p.location ILIKE $%d)", argCount, argCount+1, argCount+2)
		args = append(args, "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
		argCount += 3
	}

	if classYear != nil {
		baseQuery += fmt.Sprintf(" AND p.class_year = $%d", argCount)
		args = append(args, *classYear)
		argCount++
	}

	if house != "" {
		baseQuery += fmt.Sprintf(" AND p.house = $%d", argCount)
		args = append(args, house)
		argCount++
	}

	if location != "" {
		baseQuery += fmt.Sprintf(" AND p.location ILIKE $%d", argCount)
		args = append(args, "%"+location+"%")
		argCount++
	}

	if career != "" {
		baseQuery += fmt.Sprintf(" AND p.career ILIKE $%d", argCount)
		args = append(args, "%"+career+"%")
		argCount++
	}

	if cursor != nil {
		baseQuery += fmt.Sprintf(" AND p.created_at < (SELECT created_at FROM profiles WHERE user_id = $%d)", argCount)
		args = append(args, *cursor)
		argCount++
	}

	baseQuery += " ORDER BY p.created_at DESC"

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit+1) // Fetch one extra to check if there's more
		argCount++
	}

	query, args, err := sqlx.In(baseQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &profiles, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search profiles cursor: %w", err)
	}

	var nextCursor *string
	if len(profiles) > limit {
		profiles = profiles[:limit]
		nextCursor = &profiles[len(profiles)-1].UserID
	}

	return profiles, nextCursor, nil
}

// Delete deletes a profile
func (r *ProfileRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("profile not found")
	}
	
	return nil
}
