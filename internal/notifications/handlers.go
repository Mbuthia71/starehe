package notifications

import (
	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/repository"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationHandler(notificationRepo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
	}
}

// GetNotifications retrieves notifications for the current user
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	notifications, err := h.notificationRepo.GetNotifications(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
		"limit":         limit,
		"offset":        offset,
	})
}

// GetUnreadCount retrieves the count of unread notifications
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	count, err := h.notificationRepo.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	err := h.notificationRepo.MarkAsRead(c.Context(), notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err := h.notificationRepo.MarkAllAsRead(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	err := h.notificationRepo.DeleteNotification(c.Context(), notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification deleted",
	})
}
