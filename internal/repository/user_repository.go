package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, phone, email, password_hash, role, status, file_number)
		VALUES (:id, :phone, :email, :password_hash, :role, :status, :file_number)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, phone, email, password_hash, role, status, file_number, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByPhone retrieves a user by phone number
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	query := `SELECT id, phone, email, password_hash, role, status, file_number, created_at, updated_at FROM users WHERE phone = $1`

	err := r.db.GetContext(ctx, &user, query, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, phone, email, password_hash, role, status, file_number, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET phone = :phone, email = :email, password_hash = :password_hash, role = :role, status = :status
		WHERE id = :id
		RETURNING updated_at
	`
	
	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// UpdateStatus updates a user's status
func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, status models.UserStatus) error {
	query := `UPDATE users SET status = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// UpdateRole updates a user's role
func (r *UserRepository) UpdateRole(ctx context.Context, userID string, role models.UserRole) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `
		SELECT id, phone, email, password_hash, role, status, file_number, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`
	
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	
	return users, nil
}

// Search searches for users by name, phone, or email
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	searchQuery := `
		SELECT u.id, u.phone, u.email, u.password_hash, u.role, u.status, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE u.phone ILIKE $1 
		OR u.email ILIKE $1 
		OR p.full_name ILIKE $1
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	searchPattern := "%" + query + "%"
	err := r.db.SelectContext(ctx, &users, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	
	return users, nil
}

// Delete deletes a user (soft delete by setting status to deleted)
func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	return r.UpdateStatus(ctx, userID, models.StatusDeleted)
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	return count, nil
}

// CountByStatus returns the number of users by status
func (r *UserRepository) CountByStatus(ctx context.Context, status models.UserStatus) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE status = $1`

	err := r.db.GetContext(ctx, &count, query, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}

	return count, nil
}

// GetNextFileNumber generates the next available file number
func (r *UserRepository) GetNextFileNumber(ctx context.Context) (string, error) {
	var maxFileNumber sql.NullString
	query := `SELECT MAX(file_number) FROM users WHERE file_number IS NOT NULL`

	err := r.db.GetContext(ctx, &maxFileNumber, query)
	if err != nil {
		return "", fmt.Errorf("failed to get max file number: %w", err)
	}

	if !maxFileNumber.Valid || maxFileNumber.String == "" {
		return "00001", nil
	}

	// Parse current max and increment
	current := maxFileNumber.String
	if len(current) != 5 {
		return "00001", nil
	}

	// Convert to int, increment, and format back to 5 digits
	var num int
	_, err = fmt.Sscanf(current, "%d", &num)
	if err != nil {
		return "00001", nil
	}

	num++
	if num > 99999 {
		return "", fmt.Errorf("file number limit reached")
	}

	return fmt.Sprintf("%05d", num), nil
}
