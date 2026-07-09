package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type ConnectionRepository struct {
	db *database.DB
}

func NewConnectionRepository(db *database.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

// CreateConnection creates a new connection request
func (r *ConnectionRepository) CreateConnection(ctx context.Context, connection *models.Connection) error {
	query := `
		INSERT INTO connections (id, user_id, connected_user_id, status)
		VALUES (:id, :user_id, :connected_user_id, :status)
		RETURNING created_at, updated_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, connection)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&connection.CreatedAt, &connection.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan connection: %w", err)
		}
	}
	
	return nil
}

// GetConnection retrieves a connection between two users
func (r *ConnectionRepository) GetConnection(ctx context.Context, userID1, userID2 string) (*models.Connection, error) {
	var connection models.Connection
	query := `
		SELECT id, user_id, connected_user_id, status, created_at, updated_at
		FROM connections
		WHERE (user_id = $1 AND connected_user_id = $2)
		OR (user_id = $2 AND connected_user_id = $1)
	`
	
	err := r.db.GetContext(ctx, &connection, query, userID1, userID2)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	
	return &connection, nil
}

// GetConnectionByID retrieves a connection by ID
func (r *ConnectionRepository) GetConnectionByID(ctx context.Context, connectionID string) (*models.Connection, error) {
	var connection models.Connection
	query := `
		SELECT id, user_id, connected_user_id, status, created_at, updated_at
		FROM connections
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &connection, query, connectionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	
	return &connection, nil
}

// UpdateConnectionStatus updates a connection's status
func (r *ConnectionRepository) UpdateConnectionStatus(ctx context.Context, connectionID string, status models.ConnectionStatus) error {
	query := `UPDATE connections SET status = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, status, connectionID)
	if err != nil {
		return fmt.Errorf("failed to update connection status: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("connection not found")
	}
	
	return nil
}

// DeleteConnection deletes a connection
func (r *ConnectionRepository) DeleteConnection(ctx context.Context, connectionID string) error {
	query := `DELETE FROM connections WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, connectionID)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("connection not found")
	}
	
	return nil
}

// GetConnections retrieves connections for a user
func (r *ConnectionRepository) GetConnections(ctx context.Context, userID string, status *models.ConnectionStatus, limit, offset int) ([]*models.Connection, error) {
	var connections []*models.Connection
	
	var query string
	var args []interface{}
	
	if status != nil {
		query = `
			SELECT c.id, c.user_id, c.connected_user_id, c.status, c.created_at, c.updated_at
			FROM connections c
			WHERE (c.user_id = $1 OR c.connected_user_id = $1) AND c.status = $2
			ORDER BY c.updated_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{userID, *status, limit, offset}
	} else {
		query = `
			SELECT c.id, c.user_id, c.connected_user_id, c.status, c.created_at, c.updated_at
			FROM connections c
			WHERE c.user_id = $1 OR c.connected_user_id = $1
			ORDER BY c.updated_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, offset}
	}
	
	err := r.db.SelectContext(ctx, &connections, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}
	
	return connections, nil
}

// GetPendingConnections retrieves pending connection requests for a user
func (r *ConnectionRepository) GetPendingConnections(ctx context.Context, userID string, limit, offset int) ([]*models.Connection, error) {
	var connections []*models.Connection
	query := `
		SELECT c.id, c.user_id, c.connected_user_id, c.status, c.created_at, c.updated_at
		FROM connections c
		WHERE c.connected_user_id = $1 AND c.status = 'pending'
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &connections, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending connections: %w", err)
	}
	
	return connections, nil
}

// GetSentConnections retrieves sent connection requests for a user
func (r *ConnectionRepository) GetSentConnections(ctx context.Context, userID string, limit, offset int) ([]*models.Connection, error) {
	var connections []*models.Connection
	query := `
		SELECT c.id, c.user_id, c.connected_user_id, c.status, c.created_at, c.updated_at
		FROM connections c
		WHERE c.user_id = $1 AND c.status = 'pending'
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &connections, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent connections: %w", err)
	}
	
	return connections, nil
}

// GetConnectionCount returns the count of connections for a user
func (r *ConnectionRepository) GetConnectionCount(ctx context.Context, userID string, status *models.ConnectionStatus) (int, error) {
	var count int
	
	var query string
	var args []interface{}
	
	if status != nil {
		query = `
			SELECT COUNT(*)
			FROM connections
			WHERE (user_id = $1 OR connected_user_id = $1) AND status = $2
		`
		args = []interface{}{userID, *status}
	} else {
		query = `
			SELECT COUNT(*)
			FROM connections
			WHERE (user_id = $1 OR connected_user_id = $1)
		`
		args = []interface{}{userID}
	}
	
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection count: %w", err)
	}
	
	return count, nil
}

// Block operations
func (r *ConnectionRepository) CreateBlock(ctx context.Context, block *models.Block) error {
	query := `
		INSERT INTO blocks (id, blocker_id, blocked_id)
		VALUES (:id, :blocker_id, :blocked_id)
		RETURNING created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, block)
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&block.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan block: %w", err)
		}
	}
	
	return nil
}

func (r *ConnectionRepository) DeleteBlock(ctx context.Context, blockerID, blockedID string) error {
	query := `DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to delete block: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("block not found")
	}
	
	return nil
}

func (r *ConnectionRepository) GetBlocks(ctx context.Context, blockerID string, limit, offset int) ([]*models.Block, error) {
	var blocks []*models.Block
	query := `
		SELECT id, blocker_id, blocked_id, created_at
		FROM blocks
		WHERE blocker_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &blocks, query, blockerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}
	
	return blocks, nil
}

func (r *ConnectionRepository) IsBlocked(ctx context.Context, userID1, userID2 string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM blocks
		WHERE (blocker_id = $1 AND blocked_id = $2)
		OR (blocker_id = $2 AND blocked_id = $1)
	`
	
	err := r.db.GetContext(ctx, &count, query, userID1, userID2)
	if err != nil {
		return false, fmt.Errorf("failed to check block: %w", err)
	}
	
	return count > 0, nil
}
