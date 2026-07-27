package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/centrifugo"
	"starehian-society-platform/pkg/logger"
	"starehian-society-platform/pkg/redis"
)

type ChatService struct {
	chatRepo       *repository.ChatRepository
	connectionRepo  *repository.ConnectionRepository
	groupRepo      *repository.GroupRepository
	redis          *redis.Redis
	centrifugo     *centrifugo.CentrifugoClient
	logger         *logger.Logger
}

type CreateConversationRequest struct {
	Type string   `json:"type"`
	Name string   `json:"name,omitempty"`
	Members []string `json:"members"`
}

type SendMessageRequest struct {
	Content    *string `json:"content"`
	MediaURL   *string `json:"media_url,omitempty"`
	GroupID    *string `json:"group_id,omitempty"`
	RecipientID *string `json:"recipient_id,omitempty"`
}

func NewChatService(
	chatRepo *repository.ChatRepository,
	connectionRepo *repository.ConnectionRepository,
	groupRepo *repository.GroupRepository,
	redis *redis.Redis,
	centrifugo *centrifugo.CentrifugoClient,
	logger *logger.Logger,
) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		connectionRepo:  connectionRepo,
		groupRepo:      groupRepo,
		redis:          redis,
		centrifugo:     centrifugo,
		logger:         logger,
	}
}

// CreateDirectConversation creates a direct conversation between two users
func (s *ChatService) CreateDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error) {
	// Verify users are connected before allowing direct message
	connection, err := s.connectionRepo.GetConnection(ctx, userID1, userID2)
	if err != nil {
		s.logger.Errorf("Failed to verify connection between %s and %s: %v", userID1, userID2, err)
		return nil, fmt.Errorf("failed to verify connection: %w", err)
	}
	if connection == nil || connection.Status != "accepted" {
		return nil, fmt.Errorf("users must be connected to start a direct conversation")
	}

	// Check if conversation already exists
	existing, err := s.chatRepo.GetDirectConversation(ctx, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}
	if existing != nil {
		return existing, nil
	}

	// Create new conversation
	conversation := &models.Conversation{
		ID:   uuid.New().String(),
		Type: string(models.ConversationTypeDirect),
	}

	err = s.chatRepo.CreateConversation(ctx, conversation)
	if err != nil {
		s.logger.Errorf("Failed to create conversation: %v", err)
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Add both users as members
	member1 := &models.ConversationMember{
		ID:             uuid.New().String(),
		ConversationID: conversation.ID,
		UserID:         userID1,
		Role:           "member",
	}
	member2 := &models.ConversationMember{
		ID:             uuid.New().String(),
		ConversationID: conversation.ID,
		UserID:         userID2,
		Role:           "member",
	}

	err = s.chatRepo.AddConversationMember(ctx, member1)
	if err != nil {
		return nil, fmt.Errorf("failed to add member1: %w", err)
	}

	err = s.chatRepo.AddConversationMember(ctx, member2)
	if err != nil {
		return nil, fmt.Errorf("failed to add member2: %w", err)
	}

	s.logger.Infof("Direct conversation created: %s between %s and %s", conversation.ID, userID1, userID2)
	return conversation, nil
}

// CreateGroupConversation creates a group conversation
func (s *ChatService) CreateGroupConversation(ctx context.Context, creatorID string, req *CreateConversationRequest) (*models.Conversation, error) {
	conversation := &models.Conversation{
		ID:   uuid.New().String(),
		Type: string(models.ConversationTypeGroup),
		Name: &req.Name,
	}

	err := s.chatRepo.CreateConversation(ctx, conversation)
	if err != nil {
		s.logger.Errorf("Failed to create group conversation: %v", err)
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Add creator as admin
	adminMember := &models.ConversationMember{
		ID:             uuid.New().String(),
		ConversationID: conversation.ID,
		UserID:         creatorID,
		Role:           "admin",
	}

	err = s.chatRepo.AddConversationMember(ctx, adminMember)
	if err != nil {
		return nil, fmt.Errorf("failed to add admin: %w", err)
	}

	// Add other members if provided
	if len(req.Members) > 0 {
		for _, memberID := range req.Members {
			member := &models.ConversationMember{
				ID:             uuid.New().String(),
				ConversationID: conversation.ID,
				UserID:         memberID,
				Role:           "member",
			}

			err = s.chatRepo.AddConversationMember(ctx, member)
			if err != nil {
				s.logger.Errorf("Failed to add member %s: %v", memberID, err)
			}
		}
	}

	s.logger.Infof("Group conversation created: %s by %s", conversation.ID, creatorID)
	return conversation, nil
}

// SendMessage sends a message to a conversation or group
func (s *ChatService) SendMessage(ctx context.Context, userID, conversationID string, req *SendMessageRequest) (*models.Message, error) {
	// Handle group messages
	if req.GroupID != nil {
		return s.SendGroupMessage(ctx, userID, *req.GroupID, req)
	}

	// Handle direct messages
	if req.RecipientID != nil {
		return s.SendDirectMessage(ctx, userID, *req.RecipientID, req)
	}

	// Handle conversation messages (legacy)
	members, err := s.chatRepo.GetConversationMembers(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %w", err)
	}

	isMember := false
	for _, member := range members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, fmt.Errorf("not a member of this conversation")
	}

	if req.Content == nil && req.MediaURL == nil {
		return nil, fmt.Errorf("message must have content or media")
	}

	message := &models.Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        req.Content,
		MediaURL:       req.MediaURL,
	}

	err = s.chatRepo.CreateMessage(ctx, message)
	if err != nil {
		s.logger.Errorf("Failed to create message: %v", err)
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Broadcast via Centrifugo if available
	if s.centrifugo != nil {
		go func() {
			ctx := context.Background()
			channel := fmt.Sprintf("conversation:%s", conversationID)
			if err := s.centrifugo.Publish(ctx, channel, message); err != nil {
				s.logger.Errorf("Failed to publish to Centrifugo: %v", err)
			}
		}()
	}

	// Update sender's last read message
	s.chatRepo.UpdateLastReadMessage(ctx, conversationID, userID, message.ID)

	s.logger.Infof("Message sent: %s in conversation %s by %s", message.ID, conversationID, userID)
	return message, nil
}

// SendGroupMessage sends a message to a group
func (s *ChatService) SendGroupMessage(ctx context.Context, userID, groupID string, req *SendMessageRequest) (*models.Message, error) {
	if req.Content == nil && req.MediaURL == nil {
		return nil, fmt.Errorf("message must have content or media")
	}

	message := &models.Message{
		ID:        uuid.New().String(),
		GroupID:   &groupID,
		SenderID:  userID,
		Content:   req.Content,
		MediaURL:  req.MediaURL,
	}

	err := s.chatRepo.CreateMessage(ctx, message)
	if err != nil {
		s.logger.Errorf("Failed to create group message: %v", err)
		return nil, fmt.Errorf("failed to create group message: %w", err)
	}

	// Publish to Centrifugo channel
	channel := fmt.Sprintf("group:%s", groupID)
	if s.centrifugo != nil {
		go func() {
			ctx := context.Background()
			if err := s.centrifugo.Publish(ctx, channel, message); err != nil {
				s.logger.Errorf("Failed to publish to Centrifugo: %v", err)
			}
		}()
	}

	s.logger.Infof("Group message sent: %s in group %s by %s", message.ID, groupID, userID)
	return message, nil
}

// SendDirectMessage sends a direct message to a user
func (s *ChatService) SendDirectMessage(ctx context.Context, senderID, recipientID string, req *SendMessageRequest) (*models.Message, error) {
	if req.Content == nil && req.MediaURL == nil {
		return nil, fmt.Errorf("message must have content or media")
	}

	message := &models.Message{
		ID:          uuid.New().String(),
		RecipientID: &recipientID,
		SenderID:    senderID,
		Content:     req.Content,
		MediaURL:    req.MediaURL,
	}

	err := s.chatRepo.CreateMessage(ctx, message)
	if err != nil {
		s.logger.Errorf("Failed to create direct message: %v", err)
		return nil, fmt.Errorf("failed to create direct message: %w", err)
	}

	// Publish to Centrifugo channel for DM
	// Create sorted user IDs for consistent channel naming
	userIDs := []string{senderID, recipientID}
	if userIDs[0] > userIDs[1] {
		userIDs[0], userIDs[1] = userIDs[1], userIDs[0]
	}
	channel := fmt.Sprintf("dm:%s:%s", userIDs[0], userIDs[1])
	
	if s.centrifugo != nil {
		go func() {
			ctx := context.Background()
			if err := s.centrifugo.Publish(ctx, channel, message); err != nil {
				s.logger.Errorf("Failed to publish to Centrifugo: %v", err)
			}
		}()
	}

	s.logger.Infof("Direct message sent: %s from %s to %s", message.ID, senderID, recipientID)
	return message, nil
}

// GetMessages retrieves messages from a conversation
func (s *ChatService) GetMessages(ctx context.Context, userID, conversationID string, limit, offset int) ([]*models.Message, error) {
	// Check if user is a member
	members, err := s.chatRepo.GetConversationMembers(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %w", err)
	}

	isMember := false
	for _, member := range members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, fmt.Errorf("not a member of this conversation")
	}

	messages, err := s.chatRepo.GetMessages(ctx, conversationID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get messages: %v", err)
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

// GetConversations retrieves conversations for a user
func (s *ChatService) GetConversations(ctx context.Context, userID string, limit, offset int) ([]*models.Conversation, error) {
	conversations, err := s.chatRepo.GetConversationsByUser(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get conversations: %v", err)
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, nil
}

// MarkAsRead marks messages as read for a user
func (s *ChatService) MarkAsRead(ctx context.Context, userID, conversationID, messageID string) error {
	// Check if user is a member
	members, err := s.chatRepo.GetConversationMembers(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("failed to get conversation members: %w", err)
	}

	isMember := false
	for _, member := range members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return fmt.Errorf("not a member of this conversation")
	}

	err = s.chatRepo.UpdateLastReadMessage(ctx, conversationID, userID, messageID)
	if err != nil {
		s.logger.Errorf("Failed to update last read message: %v", err)
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	// Broadcast read receipt via Centrifugo if available
	if s.centrifugo != nil {
		go func() {
			ctx := context.Background()
			channel := fmt.Sprintf("conversation:%s", conversationID)
			readReceipt := map[string]interface{}{
				"type":      "read",
				"user_id":   userID,
				"message_id": messageID,
			}
			if err := s.centrifugo.Publish(ctx, channel, readReceipt); err != nil {
				s.logger.Errorf("Failed to publish read receipt: %v", err)
			}
		}()
	}

	s.logger.Infof("Conversation %s marked as read by %s up to message %s", conversationID, userID, messageID)
	return nil
}

// AddMember adds a member to a group conversation
func (s *ChatService) AddMember(ctx context.Context, adminID, conversationID, newMemberID string) error {
	// Check if requester is admin
	members, err := s.chatRepo.GetConversationMembers(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("failed to get conversation members: %w", err)
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == adminID && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return fmt.Errorf("only admins can add members")
	}

	// Add new member
	newMember := &models.ConversationMember{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		UserID:         newMemberID,
		Role:           "member",
	}

	err = s.chatRepo.AddConversationMember(ctx, newMember)
	if err != nil {
		s.logger.Errorf("Failed to add member: %v", err)
		return fmt.Errorf("failed to add member: %w", err)
	}

	s.logger.Infof("Member %s added to conversation %s by %s", newMemberID, conversationID, adminID)
	return nil
}

// RemoveMember removes a member from a group conversation
func (s *ChatService) RemoveMember(ctx context.Context, adminID, conversationID, memberID string) error {
	// Check if requester is admin
	members, err := s.chatRepo.GetConversationMembers(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("failed to get conversation members: %w", err)
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == adminID && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return fmt.Errorf("only admins can remove members")
	}

	err = s.chatRepo.RemoveConversationMember(ctx, conversationID, memberID)
	if err != nil {
		s.logger.Errorf("Failed to remove member: %v", err)
		return fmt.Errorf("failed to remove member: %w", err)
	}

	s.logger.Infof("Member %s removed from conversation %s by %s", memberID, conversationID, adminID)
	return nil
}

// LeaveConversation allows a user to leave a group conversation
func (s *ChatService) LeaveConversation(ctx context.Context, userID, conversationID string) error {
	err := s.chatRepo.RemoveConversationMember(ctx, conversationID, userID)
	if err != nil {
		s.logger.Errorf("Failed to leave conversation: %v", err)
		return fmt.Errorf("failed to leave conversation: %w", err)
	}

	s.logger.Infof("User %s left conversation %s", userID, conversationID)
	return nil
}

// GetGroupMessages retrieves messages from a group
func (s *ChatService) GetGroupMessages(ctx context.Context, userID, groupID string, limit, offset int) ([]*models.Message, error) {
	// Check if user is a member of the group
	if s.groupRepo != nil {
		member, err := s.groupRepo.GetGroupMember(ctx, groupID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check group membership: %w", err)
		}
		if member == nil {
			return nil, fmt.Errorf("user is not a member of this group")
		}
	}
	
	messages, err := s.chatRepo.GetGroupMessages(ctx, groupID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get group messages: %v", err)
		return nil, fmt.Errorf("failed to get group messages: %w", err)
	}

	return messages, nil
}

// GetDirectMessages retrieves direct messages between two users
func (s *ChatService) GetDirectMessages(ctx context.Context, userID1, userID2 string, limit, offset int) ([]*models.Message, error) {
	messages, err := s.chatRepo.GetDirectMessages(ctx, userID1, userID2, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get direct messages: %v", err)
		return nil, fmt.Errorf("failed to get direct messages: %w", err)
	}

	return messages, nil
}

// GetConnectionToken generates a connection token for Centrifugo
func (s *ChatService) GetConnectionToken(userID string) (string, error) {
	if s.centrifugo == nil {
		return "", fmt.Errorf("centrifugo client not configured")
	}
	return s.centrifugo.GenerateConnectionToken(userID, 0)
}

// GetChannelToken generates a channel token for a private channel
func (s *ChatService) GetChannelToken(userID, channel string) (string, error) {
	if s.centrifugo == nil {
		return "", fmt.Errorf("centrifugo client not configured")
	}
	return s.centrifugo.GenerateChannelToken(userID, channel, 0)
}
