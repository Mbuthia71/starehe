package connections

import (
	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/services"
)

type ConnectionHandler struct {
	connectionService *services.ConnectionService
}

func NewConnectionHandler(connectionService *services.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: connectionService,
	}
}

// SendConnectionRequest sends a connection request
func (h *ConnectionHandler) SendConnectionRequest(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.SendConnectionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ConnectedUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID is required",
		})
	}

	connection, err := h.connectionService.SendConnectionRequest(c.Context(), userID, req.ConnectedUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(connection)
}

// AcceptConnectionRequest accepts a connection request
func (h *ConnectionHandler) AcceptConnectionRequest(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	connectionID := c.Params("id")
	if connectionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Connection ID is required",
		})
	}

	connection, err := h.connectionService.AcceptConnectionRequest(c.Context(), userID, connectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(connection)
}

// RejectConnectionRequest rejects a connection request
func (h *ConnectionHandler) RejectConnectionRequest(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	connectionID := c.Params("id")
	if connectionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Connection ID is required",
		})
	}

	err := h.connectionService.RejectConnectionRequest(c.Context(), userID, connectionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Connection request rejected",
	})
}

// RemoveConnection removes a connection
func (h *ConnectionHandler) RemoveConnection(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID is required",
		})
	}

	err := h.connectionService.RemoveConnection(c.Context(), userID, targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Connection removed",
	})
}

// BlockUser blocks a user
func (h *ConnectionHandler) BlockUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID is required",
		})
	}

	err := h.connectionService.BlockUser(c.Context(), userID, targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User blocked",
	})
}

// UnblockUser unblocks a user
func (h *ConnectionHandler) UnblockUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	targetID := c.Params("id")
	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID is required",
		})
	}

	err := h.connectionService.UnblockUser(c.Context(), userID, targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User unblocked",
	})
}

// GetConnections retrieves connections
func (h *ConnectionHandler) GetConnections(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	statusStr := c.Query("status", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	var status *models.ConnectionStatus
	if statusStr != "" {
		s := models.ConnectionStatus(statusStr)
		status = &s
	}

	connections, err := h.connectionService.GetConnections(c.Context(), userID, status, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"connections": connections,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetPendingConnections retrieves pending connection requests
func (h *ConnectionHandler) GetPendingConnections(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	connections, err := h.connectionService.GetPendingConnections(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"connections": connections,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetSentConnections retrieves sent connection requests
func (h *ConnectionHandler) GetSentConnections(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	connections, err := h.connectionService.GetSentConnections(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"connections": connections,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetBlocks retrieves blocked users
func (h *ConnectionHandler) GetBlocks(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.Context())
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	blocks, err := h.connectionService.GetBlocks(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"blocks": blocks,
		"limit":  limit,
		"offset": offset,
	})
}
