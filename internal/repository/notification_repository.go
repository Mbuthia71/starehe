package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type NotificationRepository struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, payload)
		VALUES (:id, :user_id, :type, :payload)
		RETURNING created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&notification.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan notification: %w", err)
		}
	}
	
	return nil
}

// GetNotifications retrieves notifications for a user
func (r *NotificationRepository) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `
		SELECT id, user_id, type, payload, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &notifications, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	
	return notifications, nil
}

// GetUnreadCount retrieves the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND read_at IS NULL
	`
	
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	
	return count, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	query := `UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}
	
	return nil
}

// MarkAllAsRead marks all notifications for a user as read
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND read_at IS NULL`
	
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	
	return nil
}

// DeleteNotification deletes a notification
func (r *NotificationRepository) DeleteNotification(ctx context.Context, notificationID string) error {
	query := `DELETE FROM notifications WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}
	
	return nil
}
