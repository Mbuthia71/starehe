package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type PointsService struct {
	pointsRepo   *repository.PointsRepository
	userRepo     *repository.UserRepository
	db           *sqlx.DB
	logger       *logger.Logger
}

func NewPointsService(pointsRepo *repository.PointsRepository, userRepo *repository.UserRepository, db *sqlx.DB, appLogger *logger.Logger) *PointsService {
	return &PointsService{
		pointsRepo: pointsRepo,
		userRepo:   userRepo,
		db:         db,
		logger:     appLogger,
	}
}

// Wallet operations

func (s *PointsService) GetBalance(ctx context.Context, userID string) (*models.PointsBalanceResponse, error) {
	wallet, err := s.pointsRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &models.PointsBalanceResponse{
		Balance:             wallet.Balance,
		HeldBalance:         wallet.HeldBalance,
		Currency:            "KES",
		PointsToCurrencyRate: 0.10, // Default rate: 100 points = KES 10
		CurrencyValue:       float64(wallet.Balance) * 0.10,
	}, nil
}

func (s *PointsService) GetHistory(ctx context.Context, userID string, limit int, cursor *string) (*models.PointsHistoryResponse, error) {
	transactions, nextCursor, err := s.pointsRepo.GetUserTransactions(ctx, userID, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return &models.PointsHistoryResponse{
		Transactions: transactions,
		Cursor:       nextCursor,
		HasMore:      nextCursor != nil,
	}, nil
}

// Transaction operations

func (s *PointsService) EarnPoints(ctx context.Context, userID string, amount int64, transactionType models.TransactionType, reference *string, partnerID *string, metadata map[string]interface{}) (*models.PointsTransaction, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	wallet, err := s.pointsRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	newBalance := wallet.Balance + amount
	transactionID := uuid.New().String()
	now := time.Now()

	transaction := &models.PointsTransaction{
		ID:           transactionID,
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         transactionType,
		Amount:       amount,
		BalanceAfter: newBalance,
		Reference:    reference,
		PartnerID:    partnerID,
		Status:       models.TransactionStatusCompleted,
		Metadata:     metadata,
		CreatedAt:    now,
		ProcessedAt:  &now,
	}

	if err := s.pointsRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		return nil, err
	}

	if err := s.pointsRepo.UpdateBalance(ctx, tx, wallet.ID, newBalance, wallet.HeldBalance); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update user tier if this is an earn transaction
	if transactionType == models.TransactionTypeEarn {
		go s.updateUserTier(context.Background(), userID, amount)
	}

	return transaction, nil
}

func (s *PointsService) RedeemPoints(ctx context.Context, userID string, req *models.RedeemPointsRequest) (*models.RedeemPointsResponse, error) {
	// Validate partner
	partner, err := s.pointsRepo.GetPartner(ctx, req.PartnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}
	if partner == nil {
		return nil, fmt.Errorf("partner not found")
	}

	// Check minimum points
	if req.AmountPoints < int64(partner.MinPointsToRedeem) {
		return nil, fmt.Errorf("minimum redemption is %d points", partner.MinPointsToRedeem)
	}

	// Check daily redemption limit
	dailyCount, err := s.pointsRepo.GetUserDailyRedemptionCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check daily redemption count: %w", err)
	}
	if dailyCount >= partner.MaxDailyRedemptions {
		return nil, fmt.Errorf("daily redemption limit reached")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	wallet, err := s.pointsRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check balance
	if wallet.Balance < req.AmountPoints {
		return nil, fmt.Errorf("insufficient points balance")
	}

	// Calculate currency value
	currencyValue := float64(req.AmountPoints) * partner.PointsToCurrencyRate

	// Hold points (subtract from balance, add to held)
	newBalance := wallet.Balance - req.AmountPoints
	newHeldBalance := wallet.HeldBalance + req.AmountPoints

	transactionID := uuid.New().String()
	redemptionID := uuid.New().String()
	now := time.Now()

	// Create hold transaction
	transaction := &models.PointsTransaction{
		ID:           transactionID,
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         models.TransactionTypeHold,
		Amount:       -req.AmountPoints,
		BalanceAfter: newBalance,
		Reference:    &redemptionID,
		PartnerID:    &req.PartnerID,
		Status:       models.TransactionStatusPending,
		Metadata:     req.Metadata,
		CreatedAt:    now,
	}

	if err := s.pointsRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		return nil, err
	}

	if err := s.pointsRepo.UpdateBalance(ctx, tx, wallet.ID, newBalance, newHeldBalance); err != nil {
		return nil, err
	}

	// Create redemption record
	redemption := &models.PointsRedemption{
		ID:            redemptionID,
		TransactionID: transactionID,
		PartnerID:     &req.PartnerID,
		UserID:        userID,
		RedeemMethod:  req.Method,
		AmountPoints:  req.AmountPoints,
		CurrencyValue: &currencyValue,
		Status:        models.RedemptionStatusRequested,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.pointsRepo.CreateRedemption(ctx, tx, redemption); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Enqueue worker job for partner integration
	// This would be handled by the worker system
	// For now, return pending status

	response := &models.RedeemPointsResponse{
		RedemptionID:  redemptionID,
		TransactionID: transactionID,
		Status:        models.RedemptionStatusRequested,
		AmountPoints:  req.AmountPoints,
		CurrencyValue: currencyValue,
		Next: map[string]interface{}{
			"action": "wait_for_processing",
			"check_status": fmt.Sprintf("/api/points/redemptions/%s", redemptionID),
		},
	}

	return response, nil
}

func (s *PointsService) FinalizeRedemption(ctx context.Context, redemptionID string, status models.RedemptionStatus, partnerReference, voucherCode *string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	redemption, err := s.pointsRepo.GetRedemption(ctx, redemptionID)
	if err != nil {
		return fmt.Errorf("failed to get redemption: %w", err)
	}
	if redemption == nil {
		return fmt.Errorf("redemption not found")
	}

	transaction, err := s.pointsRepo.GetTransaction(ctx, redemption.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	wallet, err := s.pointsRepo.GetWalletByUserID(ctx, redemption.UserID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	var failureReason *string
	var voucherExpiry *time.Time

	now := time.Now()

	if status == models.RedemptionStatusConfirmed {
		// Finalize the redemption: convert hold to actual spend
		newHeldBalance := wallet.HeldBalance - redemption.AmountPoints
		
		// Update transaction to completed
		if err := s.pointsRepo.UpdateTransactionStatus(ctx, tx, transaction.ID, models.TransactionStatusCompleted, &now); err != nil {
			return err
		}

		// Update wallet held balance
		if err := s.pointsRepo.UpdateBalance(ctx, tx, wallet.ID, wallet.Balance, newHeldBalance); err != nil {
			return err
		}

		// Set voucher expiry if voucher code provided
		if voucherCode != nil {
			expiry := now.Add(30 * 24 * time.Hour) // 30 days expiry
			voucherExpiry = &expiry
		}

	} else if status == models.RedemptionStatusFailed {
		// Release held points back to balance
		newBalance := wallet.Balance + redemption.AmountPoints
		newHeldBalance := wallet.HeldBalance - redemption.AmountPoints

		// Create reversal transaction
		reversalID := uuid.New().String()
		reversal := &models.PointsTransaction{
			ID:           reversalID,
			WalletID:     wallet.ID,
			UserID:       redemption.UserID,
			Type:         models.TransactionTypeReversal,
			Amount:       redemption.AmountPoints,
			BalanceAfter: newBalance,
			Reference:    &redemptionID,
			Status:       models.TransactionStatusCompleted,
			CreatedAt:    now,
			ProcessedAt:  &now,
		}

		if err := s.pointsRepo.CreateTransaction(ctx, tx, reversal); err != nil {
			return err
		}

		// Update original transaction to failed
		if err := s.pointsRepo.UpdateTransactionStatus(ctx, tx, transaction.ID, models.TransactionStatusFailed, &now); err != nil {
			return err
		}

		// Update wallet balance
		if err := s.pointsRepo.UpdateBalance(ctx, tx, wallet.ID, newBalance, newHeldBalance); err != nil {
			return err
		}

		reason := "Partner API failed or timeout"
		failureReason = &reason
	}

	// Update redemption status
	if err := s.pointsRepo.UpdateRedemptionStatus(ctx, tx, redemptionID, status, partnerReference, failureReason); err != nil {
		return err
	}

	// Update voucher code and expiry if provided
	if voucherCode != nil {
		_, err := tx.ExecContext(ctx, `
			UPDATE points_redemptions
			SET voucher_code = $1, voucher_expiry = $2
			WHERE id = $3
		`, voucherCode, voucherExpiry, redemptionID)
		if err != nil {
			return fmt.Errorf("failed to update voucher: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Admin operations

func (s *PointsService) AdjustPoints(ctx context.Context, userID string, amount int64, reason string) (*models.PointsTransaction, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	wallet, err := s.pointsRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	newBalance := wallet.Balance + amount
	if newBalance < 0 {
		return nil, fmt.Errorf("adjustment would result in negative balance")
	}

	transactionID := uuid.New().String()
	now := time.Now()

	transaction := &models.PointsTransaction{
		ID:           transactionID,
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         models.TransactionTypeAdjust,
		Amount:       amount,
		BalanceAfter: newBalance,
		Reference:    &reason,
		Status:       models.TransactionStatusCompleted,
		Metadata:     map[string]interface{}{"admin_adjustment": true},
		CreatedAt:    now,
		ProcessedAt:  &now,
	}

	if err := s.pointsRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		return nil, err
	}

	if err := s.pointsRepo.UpdateBalance(ctx, tx, wallet.ID, newBalance, wallet.HeldBalance); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// Referral operations

func (s *PointsService) CreateReferral(ctx context.Context, userID string) (*models.ReferralResponse, error) {
	// Check if user already has an active referral
	referralCode := s.generateReferralCode()
	
	referralID := uuid.New().String()
	expiresAt := time.Now().Add(90 * 24 * time.Hour) // 90 days expiry

	referral := &models.Referral{
		ID:           referralID,
		ReferrerID:   userID,
		ReferralCode: referralCode,
		PointsAwarded: 0,
		Status:       "pending",
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now(),
	}

	if err := s.pointsRepo.CreateReferral(ctx, referral); err != nil {
		return nil, fmt.Errorf("failed to create referral: %w", err)
	}

	return &models.ReferralResponse{
		ReferralCode: referralCode,
		ReferralURL:  fmt.Sprintf("https://starehe.org/referral/%s", referralCode),
		PointsReward: 500, // 500 points for successful referral
	}, nil
}

func (s *PointsService) ProcessReferralSignup(ctx context.Context, referralCode, newUserID string) error {
	referral, err := s.pointsRepo.GetReferralByCode(ctx, referralCode)
	if err != nil {
		return fmt.Errorf("failed to get referral: %w", err)
	}
	if referral == nil {
		return fmt.Errorf("invalid referral code")
	}
	if referral.Status != "pending" {
		return fmt.Errorf("referral already used or expired")
	}
	if referral.ExpiresAt != nil && time.Now().After(*referral.ExpiresAt) {
		return fmt.Errorf("referral code expired")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Award points to referrer
	pointsReward := int64(500)
	reference := fmt.Sprintf("referral:%s", referralCode)
	
	_, err = s.EarnPoints(ctx, referral.ReferrerID, pointsReward, models.TransactionTypeEarn, &reference, nil, map[string]interface{}{
		"type": "referral",
		"referred_user_id": newUserID,
	})
	if err != nil {
		return fmt.Errorf("failed to award referral points: %w", err)
	}

	// Award points to new user
	newUserReference := fmt.Sprintf("referral_signup:%s", referralCode)
	_, err = s.EarnPoints(ctx, newUserID, pointsReward, models.TransactionTypeEarn, &newUserReference, nil, map[string]interface{}{
		"type": "referral_signup",
		"referrer_id": referral.ReferrerID,
	})
	if err != nil {
		return fmt.Errorf("failed to award signup points: %w", err)
	}

	// Update referral record
	now := time.Now()
	if err := s.pointsRepo.UpdateReferral(ctx, tx, referral.ID, newUserID, int(pointsReward), "completed", &now); err != nil {
		return err
	}

	// Award super referrer badge if user has 5+ successful referrals
	go s.checkSuperReferrerBadge(context.Background(), referral.ReferrerID)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Gamification operations

func (s *PointsService) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	return s.pointsRepo.GetUserBadges(ctx, userID)
}

func (s *PointsService) GetUserTier(ctx context.Context, userID string) (*models.UserTierRecord, error) {
	return s.pointsRepo.GetUserTier(ctx, userID)
}

func (s *PointsService) updateUserTier(ctx context.Context, userID string, pointsEarned int64) {
	tier, err := s.pointsRepo.GetUserTier(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user tier: %v", err)
		return
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Errorf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	newLifetimePoints := pointsEarned
	newTier := models.UserTierBronze
	progress := 0

	if tier != nil {
		newLifetimePoints = tier.PointsEarnedLifetime + pointsEarned
	}

	// Calculate tier based on lifetime points
	switch {
	case newLifetimePoints >= 50000:
		newTier = models.UserTierDiamond
		progress = int(newLifetimePoints % 10000)
	case newLifetimePoints >= 25000:
		newTier = models.UserTierPlatinum
		progress = int(newLifetimePoints % 25000)
	case newLifetimePoints >= 10000:
		newTier = models.UserTierGold
		progress = int(newLifetimePoints % 10000)
	case newLifetimePoints >= 5000:
		newTier = models.UserTierSilver
		progress = int(newLifetimePoints % 5000)
	default:
		newTier = models.UserTierBronze
		progress = int(newLifetimePoints)
	}

	tierRecord := &models.UserTierRecord{
		ID:                 uuid.New().String(),
		UserID:             userID,
		Tier:               newTier,
		PointsEarnedLifetime: newLifetimePoints,
		CurrentLevelProgress: progress,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if tier != nil && tier.Tier != newTier {
		now := time.Now()
		tierRecord.PromotedAt = &now
	}

	if err := s.pointsRepo.CreateOrUpdateUserTier(ctx, tx, tierRecord); err != nil {
		s.logger.Errorf("Failed to update user tier: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("Failed to commit tier update: %v", err)
		return
	}
}

func (s *PointsService) checkSuperReferrerBadge(ctx context.Context, userID string) {
	// Count successful referrals
	// This would require a query to count referrals with status 'completed'
	// For now, this is a placeholder
}

// Campaign operations

func (s *PointsService) GetActiveCampaigns(ctx context.Context) ([]models.PointsCampaign, error) {
	return s.pointsRepo.GetActiveCampaigns(ctx)
}

func (s *PointsService) AwardCampaignPoints(ctx context.Context, userID, campaignID string) error {
	campaigns, err := s.pointsRepo.GetActiveCampaigns(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active campaigns: %w", err)
	}

	var campaign *models.PointsCampaign
	for _, c := range campaigns {
		if c.ID == campaignID {
			campaign = &c
			break
		}
	}

	if campaign == nil {
		return fmt.Errorf("campaign not found or inactive")
	}

	// Check if user already awarded
	awardCount, err := s.pointsRepo.GetCampaignAwardCount(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to check award count: %w", err)
	}

	if campaign.MaxAwards != nil && awardCount >= *campaign.MaxAwards {
		return fmt.Errorf("campaign award limit reached")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Award points
	reference := fmt.Sprintf("campaign:%s", campaignID)
	transaction, err := s.EarnPoints(ctx, userID, int64(campaign.PointsReward), models.TransactionTypeEarn, &reference, nil, map[string]interface{}{
		"type": "campaign",
		"campaign_id": campaignID,
	})
	if err != nil {
		return fmt.Errorf("failed to award campaign points: %w", err)
	}

	// Record award
	award := &models.CampaignAward{
		ID:            uuid.New().String(),
		CampaignID:    campaignID,
		UserID:        userID,
		TransactionID: &transaction.ID,
		AwardedAt:     time.Now(),
	}

	if err := s.pointsRepo.CreateCampaignAward(ctx, tx, award); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper functions

func (s *PointsService) generateReferralCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
