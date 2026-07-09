package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type ConnectionService struct {
	connectionRepo *repository.ConnectionRepository
	userRepo       *repository.UserRepository
	logger         *logger.Logger
}

type SendConnectionRequest struct {
	ConnectedUserID string `json:"connected_user_id"`
}

type UpdateConnectionRequest struct {
	Status models.ConnectionStatus `json:"status"`
}

func NewConnectionService(
	connectionRepo *repository.ConnectionRepository,
	userRepo *repository.UserRepository,
	logger *logger.Logger,
) *ConnectionService {
	return &ConnectionService{
		connectionRepo: connectionRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

// SendConnectionRequest sends a connection request
func (s *ConnectionService) SendConnectionRequest(ctx context.Context, userID, targetID string) (*models.Connection, error) {
	// Check if users exist
	targetUser, err := s.userRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target user: %w", err)
	}
	if targetUser == nil {
		return nil, fmt.Errorf("target user not found")
	}

	// Check if can send request
	canSend, err := s.authzService.CanSendConnectionRequest(ctx, userID, targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}
	if !canSend {
		return nil, fmt.Errorf("cannot send connection request")
	}

	// Check if connection already exists
	existing, err := s.connectionRepo.GetConnection(ctx, userID, targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing connection: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("connection already exists")
	}

	// Create connection request
	connection := &models.Connection{
		ID:              uuid.New().String(),
		UserID:          userID,
		ConnectedUserID: targetID,
		Status:          models.ConnectionStatusPending,
	}

	err = s.connectionRepo.CreateConnection(ctx, connection)
	if err != nil {
		s.logger.Errorf("Failed to create connection: %v", err)
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	s.logger.Infof("Connection request sent: %s -> %s", userID, targetID)
	return connection, nil
}

// AcceptConnectionRequest accepts a connection request
func (s *ConnectionService) AcceptConnectionRequest(ctx context.Context, userID, connectionID string) (*models.Connection, error) {
	// Get connection
	connection, err := s.connectionRepo.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	if connection == nil {
		return nil, fmt.Errorf("connection not found")
	}

	// Check if user is the recipient
	if connection.ConnectedUserID != userID {
		return nil, fmt.Errorf("not authorized to accept this request")
	}

	// Update status
	err = s.connectionRepo.UpdateConnectionStatus(ctx, connectionID, models.ConnectionStatusAccepted)
	if err != nil {
		s.logger.Errorf("Failed to update connection status: %v", err)
		return nil, fmt.Errorf("failed to accept connection: %w", err)
	}

	connection.Status = models.ConnectionStatusAccepted
	s.logger.Infof("Connection accepted: %s", connectionID)
	return connection, nil
}

// RejectConnectionRequest rejects a connection request
func (s *ConnectionService) RejectConnectionRequest(ctx context.Context, userID, connectionID string) error {
	// Get connection
	connection, err := s.connectionRepo.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	if connection == nil {
		return fmt.Errorf("connection not found")
	}

	// Check if user is the recipient
	if connection.ConnectedUserID != userID {
		return fmt.Errorf("not authorized to reject this request")
	}

	// Delete connection
	err = s.connectionRepo.DeleteConnection(ctx, connectionID)
	if err != nil {
		s.logger.Errorf("Failed to delete connection: %v", err)
		return fmt.Errorf("failed to reject connection: %w", err)
	}

	s.logger.Infof("Connection rejected: %s", connectionID)
	return nil
}

// RemoveConnection removes a connection
func (s *ConnectionService) RemoveConnection(ctx context.Context, userID, targetID string) error {
	// Get connection
	connection, err := s.connectionRepo.GetConnection(ctx, userID, targetID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	if connection == nil {
		return fmt.Errorf("connection not found")
	}

	// Check if user is part of the connection
	if connection.UserID != userID && connection.ConnectedUserID != userID {
		return fmt.Errorf("not authorized to remove this connection")
	}

	// Delete connection
	err = s.connectionRepo.DeleteConnection(ctx, connection.ID)
	if err != nil {
		s.logger.Errorf("Failed to delete connection: %v", err)
		return fmt.Errorf("failed to remove connection: %w", err)
	}

	s.logger.Infof("Connection removed: %s", connection.ID)
	return nil
}

// BlockUser blocks a user
func (s *ConnectionService) BlockUser(ctx context.Context, blockerID, blockedID string) error {
	if blockerID == blockedID {
		return fmt.Errorf("cannot block yourself")
	}

	// Check if already blocked
	isBlocked, err := s.connectionRepo.IsBlocked(ctx, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to check block: %w", err)
	}
	if isBlocked {
		return fmt.Errorf("user already blocked")
	}

	// Create block
	block := &models.Block{
		ID:        uuid.New().String(),
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	err = s.connectionRepo.CreateBlock(ctx, block)
	if err != nil {
		s.logger.Errorf("Failed to create block: %v", err)
		return fmt.Errorf("failed to block user: %w", err)
	}

	// Remove any existing connection
	connection, _ := s.connectionRepo.GetConnection(ctx, blockerID, blockedID)
	if connection != nil {
		s.connectionRepo.DeleteConnection(ctx, connection.ID)
	}

	s.logger.Infof("User blocked: %s -> %s", blockerID, blockedID)
	return nil
}

// UnblockUser unblocks a user
func (s *ConnectionService) UnblockUser(ctx context.Context, blockerID, blockedID string) error {
	err := s.connectionRepo.DeleteBlock(ctx, blockerID, blockedID)
	if err != nil {
		s.logger.Errorf("Failed to delete block: %v", err)
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	s.logger.Infof("User unblocked: %s -> %s", blockerID, blockedID)
	return nil
}

// GetConnections retrieves connections for a user
func (s *ConnectionService) GetConnections(ctx context.Context, userID string, status *models.ConnectionStatus, limit, offset int) ([]*models.Connection, error) {
	connections, err := s.connectionRepo.GetConnections(ctx, userID, status, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get connections: %v", err)
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	return connections, nil
}

// GetPendingConnections retrieves pending connection requests
func (s *ConnectionService) GetPendingConnections(ctx context.Context, userID string, limit, offset int) ([]*models.Connection, error) {
	connections, err := s.connectionRepo.GetPendingConnections(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get pending connections: %v", err)
		return nil, fmt.Errorf("failed to get pending connections: %w", err)
	}

	return connections, nil
}

// GetSentConnections retrieves sent connection requests
func (s *ConnectionService) GetSentConnections(ctx context.Context, userID string, limit, offset int) ([]*models.Connection, error) {
	connections, err := s.connectionRepo.GetSentConnections(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get sent connections: %v", err)
		return nil, fmt.Errorf("failed to get sent connections: %w", err)
	}

	return connections, nil
}

// GetBlocks retrieves blocked users
func (s *ConnectionService) GetBlocks(ctx context.Context, userID string, limit, offset int) ([]*models.Block, error) {
	blocks, err := s.connectionRepo.GetBlocks(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get blocks: %v", err)
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}

	return blocks, nil
}
