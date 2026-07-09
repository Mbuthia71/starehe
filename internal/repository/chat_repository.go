package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type ChatRepository struct {
	db *database.DB
}

func NewChatRepository(db *database.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Conversations
func (r *ChatRepository) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	query := `
		INSERT INTO conversations (id, type, name)
		VALUES (:id, :type, :name)
		RETURNING created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, conversation)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&conversation.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan conversation: %w", err)
		}
	}
	
	return nil
}

func (r *ChatRepository) GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error) {
	var conversation models.Conversation
	query := `
		SELECT id, type, name, created_at
		FROM conversations
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &conversation, query, conversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	
	return &conversation, nil
}

func (r *ChatRepository) GetDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error) {
	var conversation models.Conversation
	query := `
		SELECT c.id, c.type, c.name, c.created_at
		FROM conversations c
		WHERE c.type = 'direct'
		AND EXISTS (
			SELECT 1 FROM conversation_members cm1
			WHERE cm1.conversation_id = c.id AND cm1.user_id = $1
		)
		AND EXISTS (
			SELECT 1 FROM conversation_members cm2
			WHERE cm2.conversation_id = c.id AND cm2.user_id = $2
		)
		AND (
			SELECT COUNT(*) FROM conversation_members cm
			WHERE cm.conversation_id = c.id
		) = 2
		LIMIT 1
	`
	
	err := r.db.GetContext(ctx, &conversation, query, userID1, userID2)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get direct conversation: %w", err)
	}
	
	return &conversation, nil
}

func (r *ChatRepository) GetConversationsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	query := `
		SELECT DISTINCT c.id, c.type, c.name, c.created_at
		FROM conversations c
		INNER JOIN conversation_members cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &conversations, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	
	return conversations, nil
}

// Conversation Members
func (r *ChatRepository) AddConversationMember(ctx context.Context, member *models.ConversationMember) error {
	query := `
		INSERT INTO conversation_members (id, conversation_id, user_id, role)
		VALUES (:id, :conversation_id, :user_id, :role)
		RETURNING joined_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, member)
	if err != nil {
		return fmt.Errorf("failed to add conversation member: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&member.JoinedAt); err != nil {
			return fmt.Errorf("failed to scan conversation member: %w", err)
		}
	}
	
	return nil
}

func (r *ChatRepository) GetConversationMembers(ctx context.Context, conversationID string) ([]*models.ConversationMember, error) {
	var members []*models.ConversationMember
	query := `
		SELECT id, conversation_id, user_id, role, joined_at, last_read_message_id
		FROM conversation_members
		WHERE conversation_id = $1
		ORDER BY joined_at ASC
	`
	
	err := r.db.SelectContext(ctx, &members, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation members: %w", err)
	}
	
	return members, nil
}

func (r *ChatRepository) UpdateLastReadMessage(ctx context.Context, conversationID, userID, messageID string) error {
	query := `
		UPDATE conversation_members
		SET last_read_message_id = $1
		WHERE conversation_id = $2 AND user_id = $3
	`
	
	result, err := r.db.ExecContext(ctx, query, messageID, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to update last read message: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("conversation member not found")
	}
	
	return nil
}

func (r *ChatRepository) RemoveConversationMember(ctx context.Context, conversationID, userID string) error {
	query := `DELETE FROM conversation_members WHERE conversation_id = $1 AND user_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove conversation member: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("conversation member not found")
	}
	
	return nil
}

// Messages
func (r *ChatRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, sender_id, content, media_url)
		VALUES (:id, :conversation_id, :sender_id, :content, :media_url)
		RETURNING created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&message.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan message: %w", err)
		}
	}
	
	return nil
}

func (r *ChatRepository) GetMessage(ctx context.Context, messageID string) (*models.Message, error) {
	var message models.Message
	query := `
		SELECT id, conversation_id, sender_id, content, media_url, created_at
		FROM messages
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &message, query, messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	
	return &message, nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	query := `
		SELECT id, conversation_id, sender_id, content, media_url, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &messages, query, conversationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	
	return messages, nil
}

func (r *ChatRepository) GetMessagesAfter(ctx context.Context, conversationID, afterMessageID string, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	query := `
		SELECT id, conversation_id, sender_id, content, media_url, created_at
		FROM messages
		WHERE conversation_id = $1 AND created_at > (SELECT created_at FROM messages WHERE id = $2)
		ORDER BY created_at ASC
		LIMIT $3
	`
	
	err := r.db.SelectContext(ctx, &messages, query, conversationID, afterMessageID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages after: %w", err)
	}
	
	return messages, nil
}

func (r *ChatRepository) DeleteMessage(ctx context.Context, messageID string) error {
	query := `DELETE FROM messages WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("message not found")
	}
	
	return nil
}
