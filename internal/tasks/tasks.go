package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/internal/services"
	"starehian-society-platform/pkg/logger"
)

// Task type constants
const (
	TypeProcessRedemption       = "points:process_redemption"
	TypeProcessPartnerCallback  = "points:process_partner_callback"
	TypeSendEmail              = "email:send"
	TypeSendPushNotification   = "notification:push"
	TypeProcessMedia           = "media:process"
	TypeAggregateAnalytics     = "analytics:aggregate"
)

// TaskHandler handles background tasks
type TaskHandler struct {
	pointsService     *services.PointsService
	notificationRepo  *repository.NotificationRepository
	logger            *logger.Logger
}

func NewTaskHandler(pointsService *services.PointsService, notificationRepo *repository.NotificationRepository, appLogger *logger.Logger) *TaskHandler {
	return &TaskHandler{
		pointsService:    pointsService,
		notificationRepo: notificationRepo,
		logger:          appLogger,
	}
}

// ProcessRedemption handles points redemption processing
func (h *TaskHandler) ProcessRedemption(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		RedemptionID string `json:"redemption_id"`
		PartnerID    string `json:"partner_id"`
		UserID       string `json:"user_id"`
		AmountPoints int64  `json:"amount_points"`
		Method       string `json:"method"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Processing redemption %s for user %s", payload.RedemptionID, payload.UserID)

	// Simulate partner API call
	// In production, this would call the actual partner API (e.g., Naivas)
	success := h.callPartnerAPI(payload.PartnerID, payload.AmountPoints, payload.Method)

	var status models.RedemptionStatus
	var partnerReference *string
	var voucherCode *string

	if success {
		status = models.RedemptionStatusConfirmed
		ref := fmt.Sprintf("PARTNER-%s", payload.RedemptionID[:8])
		partnerReference = &ref

		if payload.Method == "voucher" {
			code := h.generateVoucherCode()
			voucherCode = &code
		}
	} else {
		status = models.RedemptionStatusFailed
		reason := "Partner API timeout or error"
		partnerReference = &reason
	}

	if err := h.pointsService.FinalizeRedemption(ctx, payload.RedemptionID, status, partnerReference, voucherCode); err != nil {
		h.logger.Errorf("Failed to finalize redemption %s: %v", payload.RedemptionID, err)
		return err
	}

	h.logger.Infof("Redemption %s processed with status %s", payload.RedemptionID, status)

	// Send notification to user
	if status == models.RedemptionStatusConfirmed {
		h.sendRedemptionNotification(ctx, payload.UserID, payload.AmountPoints, voucherCode)
	}

	return nil
}

// ProcessPartnerCallback handles callbacks from partner systems
func (h *TaskHandler) ProcessPartnerCallback(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		RedemptionID string `json:"redemption_id"`
		Status       string `json:"status"`
		Reference    string `json:"reference"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Processing partner callback for redemption %s", payload.RedemptionID)

	var status models.RedemptionStatus
	switch payload.Status {
	case "confirmed":
		status = models.RedemptionStatusConfirmed
	case "failed":
		status = models.RedemptionStatusFailed
	default:
		status = models.RedemptionStatusFailed
	}

	reference := &payload.Reference
	if err := h.pointsService.FinalizeRedemption(ctx, payload.RedemptionID, status, reference, nil); err != nil {
		h.logger.Errorf("Failed to process callback for redemption %s: %v", payload.RedemptionID, err)
		return err
	}

	return nil
}

// SendEmail handles email sending tasks
func (h *TaskHandler) SendEmail(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Sending email to %s: %s", payload.To, payload.Subject)

	// In production, integrate with email service (e.g., SendGrid, AWS SES)
	// For now, just log
	time.Sleep(100 * time.Millisecond) // Simulate email sending

	return nil
}

// SendPushNotification handles push notification tasks
func (h *TaskHandler) SendPushNotification(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		UserID  string                 `json:"user_id"`
		Title   string                 `json:"title"`
		Body    string                 `json:"body"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Sending push notification to user %s: %s", payload.UserID, payload.Title)

	// In production, integrate with FCM or APNs
	// For now, just log
	time.Sleep(50 * time.Millisecond) // Simulate push notification

	return nil
}

// ProcessMedia handles media processing tasks (thumbnailing, transcoding)
func (h *TaskHandler) ProcessMedia(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		MediaURL string `json:"media_url"`
		Type     string `json:"type"` // image, video
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Processing media %s of type %s", payload.MediaURL, payload.Type)

	// In production, integrate with media processing service
	// For now, just log
	time.Sleep(500 * time.Millisecond) // Simulate media processing

	return nil
}

// AggregateAnalytics handles analytics aggregation tasks
func (h *TaskHandler) AggregateAnalytics(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		Date string `json:"date"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	h.logger.Infof("Aggregating analytics for date %s", payload.Date)

	// In production, aggregate metrics from various sources
	// For now, just log
	time.Sleep(1 * time.Second) // Simulate analytics aggregation

	return nil
}

// Helper functions

func (h *TaskHandler) callPartnerAPI(partnerID string, amount int64, method string) bool {
	// Simulate partner API call
	// In production, this would:
	// 1. Get partner details from DB
	// 2. Sign request with partner's API key
	// 3. Call partner's API endpoint
	// 4. Handle response and retry on failure
	
	time.Sleep(200 * time.Millisecond) // Simulate API call
	
	// Simulate 90% success rate
	return time.Now().UnixNano()%10 < 9
}

func (h *TaskHandler) generateVoucherCode() string {
	// Generate a human-readable voucher code
	// Format: XXXX-XXXX-XXXX
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return fmt.Sprintf("%s-%s-%s", string(b[0:4]), string(b[4:8]), string(b[8:12]))
}

func (h *TaskHandler) sendRedemptionNotification(ctx context.Context, userID string, amount int64, voucherCode *string) {
	payload := map[string]interface{}{
		"user_id": userID,
		"title":   "Points Redeemed Successfully",
		"body":    fmt.Sprintf("You have successfully redeemed %d points!", amount),
		"data": map[string]interface{}{
			"type":  "redemption",
			"amount": amount,
		},
	}

	if voucherCode != nil {
		payload["data"].(map[string]interface{})["voucher_code"] = *voucherCode
	}

	task, err := asynq.NewTask(TypeSendPushNotification, payload)
	if err != nil {
		h.logger.Errorf("Failed to create notification task: %v", err)
		return
	}

	// In production, enqueue the task
	// For now, just log
	h.logger.Infof("Would enqueue notification task for user %s", userID)
}
