package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"starehian-society-platform/internal/models"
)

type PointsRepository struct {
	db *sqlx.DB
}

func NewPointsRepository(db *sqlx.DB) *PointsRepository {
	return &PointsRepository{db: db}
}

// Wallet operations

func (r *PointsRepository) GetWalletByUserID(ctx context.Context, userID string) (*models.PointsWallet, error) {
	var wallet models.PointsWallet
	err := r.db.GetContext(ctx, &wallet, `
		SELECT id, user_id, balance, held_balance, created_at, updated_at
		FROM points_wallets
		WHERE user_id = $1
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return &wallet, nil
}

func (r *PointsRepository) CreateWallet(ctx context.Context, userID string) (*models.PointsWallet, error) {
	walletID := uuid.New().String()
	now := time.Now()
	
	wallet := &models.PointsWallet{
		ID:          walletID,
		UserID:      userID,
		Balance:     0,
		HeldBalance: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO points_wallets (id, user_id, balance, held_balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, wallet.ID, wallet.UserID, wallet.Balance, wallet.HeldBalance, wallet.CreatedAt, wallet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (r *PointsRepository) GetOrCreateWallet(ctx context.Context, userID string) (*models.PointsWallet, error) {
	wallet, err := r.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		return wallet, nil
	}
	return r.CreateWallet(ctx, userID)
}

func (r *PointsRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, walletID string, newBalance, newHeldBalance int64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE points_wallets
		SET balance = $1, held_balance = $2, updated_at = NOW()
		WHERE id = $3
	`, newBalance, newHeldBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}

// Transaction operations

func (r *PointsRepository) CreateTransaction(ctx context.Context, tx *sqlx.Tx, transaction *models.PointsTransaction) error {
	metadataJSON, err := json.Marshal(transaction.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction metadata: %w", err)
	}
	
	query := `
		INSERT INTO points_transactions 
		(id, wallet_id, user_id, type, amount, balance_after, reference, partner_id, status, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err = tx.ExecContext(ctx, query,
		transaction.ID,
		transaction.WalletID,
		transaction.UserID,
		transaction.Type,
		transaction.Amount,
		transaction.BalanceAfter,
		transaction.Reference,
		transaction.PartnerID,
		transaction.Status,
		metadataJSON,
		transaction.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetTransaction(ctx context.Context, transactionID string) (*models.PointsTransaction, error) {
	var transaction models.PointsTransaction
	var metadataJSON []byte
	
	err := r.db.GetContext(ctx, &transaction, `
		SELECT id, wallet_id, user_id, type, amount, balance_after, reference, partner_id, status, metadata, created_at, processed_at
		FROM points_transactions
		WHERE id = $1
	`, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	
	if metadataJSON != nil {
		if err := json.Unmarshal(metadataJSON, &transaction.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction metadata: %w", err)
		}
	}
	
	return &transaction, nil
}

func (r *PointsRepository) UpdateTransactionStatus(ctx context.Context, tx *sqlx.Tx, transactionID string, status models.TransactionStatus, processedAt *time.Time) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE points_transactions
		SET status = $1, processed_at = $2
		WHERE id = $3
	`, status, processedAt, transactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetUserTransactions(ctx context.Context, userID string, limit int, cursor *string) ([]models.PointsTransaction, *string, error) {
	query := `
		SELECT id, wallet_id, user_id, type, amount, balance_after, reference, partner_id, status, metadata, created_at, processed_at
		FROM points_transactions
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	argCount := 1
	
	if cursor != nil {
		query += fmt.Sprintf(" AND created_at < (SELECT created_at FROM points_transactions WHERE id = $%d)", argCount+1)
		args = append(args, *cursor)
		argCount++
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit+1) // Fetch one extra to check if there's more
		argCount++
	}
	
	var transactions []models.PointsTransaction
	err := r.db.SelectContext(ctx, &transactions, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	
	var nextCursor *string
	if len(transactions) > limit {
		transactions = transactions[:limit]
		nextCursor = &transactions[len(transactions)-1].ID
	}
	
	return transactions, nextCursor, nil
}

// Partner operations

func (r *PointsRepository) CreatePartner(ctx context.Context, partner *models.Partner) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO partners 
		(id, name, description, integration_mode, api_base_url, api_key_encrypted, callback_secret, logo_url, active, points_to_currency_rate, min_points_to_redeem, max_daily_redemptions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, partner.ID, partner.Name, partner.Description, partner.IntegrationMode, partner.APIBaseURL, 
		partner.APIKeyEncrypted, partner.CallbackSecret, partner.LogoURL, partner.Active, 
		partner.PointsToCurrencyRate, partner.MinPointsToRedeem, partner.MaxDailyRedemptions, 
		partner.CreatedAt, partner.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create partner: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetPartner(ctx context.Context, partnerID string) (*models.Partner, error) {
	var partner models.Partner
	err := r.db.GetContext(ctx, &partner, `
		SELECT id, name, description, integration_mode, api_base_url, api_key_encrypted, callback_secret, logo_url, active, points_to_currency_rate, min_points_to_redeem, max_daily_redemptions, created_at, updated_at
		FROM partners
		WHERE id = $1 AND active = true
	`, partnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}
	return &partner, nil
}

func (r *PointsRepository) ListActivePartners(ctx context.Context) ([]models.Partner, error) {
	var partners []models.Partner
	err := r.db.SelectContext(ctx, &partners, `
		SELECT id, name, description, integration_mode, api_base_url, logo_url, active, points_to_currency_rate, min_points_to_redeem, max_daily_redemptions, created_at, updated_at
		FROM partners
		WHERE active = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list partners: %w", err)
	}
	return partners, nil
}

// Redemption operations

func (r *PointsRepository) CreateRedemption(ctx context.Context, tx *sqlx.Tx, redemption *models.PointsRedemption) error {
	query := `
		INSERT INTO points_redemptions 
		(id, transaction_id, partner_id, user_id, redeem_method, amount_points, currency_value, partner_reference, voucher_code, voucher_expiry, status, failure_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	
	_, err := tx.ExecContext(ctx, query,
		redemption.ID, redemption.TransactionID, redemption.PartnerID, redemption.UserID,
		redemption.RedeemMethod, redemption.AmountPoints, redemption.CurrencyValue,
		redemption.PartnerReference, redemption.VoucherCode, redemption.VoucherExpiry,
		redemption.Status, redemption.FailureReason, redemption.CreatedAt, redemption.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create redemption: %w", err)
	}
	return nil
}

func (r *PointsRepository) UpdateRedemptionStatus(ctx context.Context, tx *sqlx.Tx, redemptionID string, status models.RedemptionStatus, partnerReference, failureReason *string) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE points_redemptions
		SET status = $1, partner_reference = $2, failure_reason = $3, updated_at = NOW()
		WHERE id = $4
	`, status, partnerReference, failureReason, redemptionID)
	if err != nil {
		return fmt.Errorf("failed to update redemption status: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetRedemption(ctx context.Context, redemptionID string) (*models.PointsRedemption, error) {
	var redemption models.PointsRedemption
	err := r.db.GetContext(ctx, &redemption, `
		SELECT id, transaction_id, partner_id, user_id, redeem_method, amount_points, currency_value, partner_reference, voucher_code, voucher_expiry, status, failure_reason, created_at, updated_at
		FROM points_redemptions
		WHERE id = $1
	`, redemptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get redemption: %w", err)
	}
	return &redemption, nil
}

func (r *PointsRepository) GetUserRedemptions(ctx context.Context, userID string, limit int, cursor *string) ([]models.PointsRedemption, *string, error) {
	query := `
		SELECT id, transaction_id, partner_id, user_id, redeem_method, amount_points, currency_value, partner_reference, voucher_code, voucher_expiry, status, failure_reason, created_at, updated_at
		FROM points_redemptions
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	argCount := 1
	
	if cursor != nil {
		query += fmt.Sprintf(" AND created_at < (SELECT created_at FROM points_redemptions WHERE id = $%d)", argCount+1)
		args = append(args, *cursor)
		argCount++
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit+1)
		argCount++
	}
	
	var redemptions []models.PointsRedemption
	err := r.db.SelectContext(ctx, &redemptions, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user redemptions: %w", err)
	}
	
	var nextCursor *string
	if len(redemptions) > limit {
		redemptions = redemptions[:limit]
		nextCursor = &redemptions[len(redemptions)-1].ID
	}
	
	return redemptions, nextCursor, nil
}

func (r *PointsRepository) GetUserDailyRedemptionCount(ctx context.Context, userID string) (int, error) {
	var count int
	today := time.Now().Truncate(24 * time.Hour)
	
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		FROM points_redemptions
		WHERE user_id = $1 AND created_at >= $2 AND status IN ('requested', 'reserved', 'sent_to_partner', 'confirmed')
	`, userID, today)
	if err != nil {
		return 0, fmt.Errorf("failed to get daily redemption count: %w", err)
	}
	return count, nil
}

// Referral operations

func (r *PointsRepository) CreateReferral(ctx context.Context, referral *models.Referral) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO referrals (id, referrer_id, referral_code, status, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, referral.ID, referral.ReferrerID, referral.ReferralCode, referral.Status, referral.ExpiresAt, referral.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetReferralByCode(ctx context.Context, code string) (*models.Referral, error) {
	var referral models.Referral
	err := r.db.GetContext(ctx, &referral, `
		SELECT id, referrer_id, referred_id, referral_code, points_awarded, status, completed_at, expires_at, created_at
		FROM referrals
		WHERE referral_code = $1
	`, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get referral by code: %w", err)
	}
	return &referral, nil
}

func (r *PointsRepository) UpdateReferral(ctx context.Context, tx *sqlx.Tx, referralID, referredID string, pointsAwarded int, status string, completedAt *time.Time) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE referrals
		SET referred_id = $1, points_awarded = $2, status = $3, completed_at = $4
		WHERE id = $5
	`, referredID, pointsAwarded, status, completedAt, referralID)
	if err != nil {
		return fmt.Errorf("failed to update referral: %w", err)
	}
	return nil
}

// Gamification operations

func (r *PointsRepository) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	err := r.db.SelectContext(ctx, &badges, `
		SELECT id, user_id, badge_type, earned_at, metadata
		FROM user_badges
		WHERE user_id = $1
		ORDER BY earned_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user badges: %w", err)
	}
	return badges, nil
}

func (r *PointsRepository) AwardBadge(ctx context.Context, tx *sqlx.Tx, badge *models.UserBadge) error {
	metadataJSON, err := json.Marshal(badge.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal badge metadata: %w", err)
	}
	
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_badges (id, user_id, badge_type, earned_at, metadata)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, badge_type) DO NOTHING
	`, badge.ID, badge.UserID, badge.BadgeType, badge.EarnedAt, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to award badge: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetUserTier(ctx context.Context, userID string) (*models.UserTierRecord, error) {
	var tier models.UserTierRecord
	err := r.db.GetContext(ctx, &tier, `
		SELECT id, user_id, tier, points_earned_lifetime, current_level_progress, promoted_at, created_at, updated_at
		FROM user_tiers
		WHERE user_id = $1
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user tier: %w", err)
	}
	return &tier, nil
}

func (r *PointsRepository) CreateOrUpdateUserTier(ctx context.Context, tx *sqlx.Tx, tier *models.UserTierRecord) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO user_tiers (id, user_id, tier, points_earned_lifetime, current_level_progress, promoted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			tier = excluded.tier,
			points_earned_lifetime = excluded.points_earned_lifetime,
			current_level_progress = excluded.current_level_progress,
			promoted_at = excluded.promoted_at,
			updated_at = NOW()
	`, tier.ID, tier.UserID, tier.Tier, tier.PointsEarnedLifetime, tier.CurrentLevelProgress, tier.PromotedAt, tier.CreatedAt, tier.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create or update user tier: %w", err)
	}
	return nil
}

// Campaign operations

func (r *PointsRepository) GetActiveCampaigns(ctx context.Context) ([]models.PointsCampaign, error) {
	var campaigns []models.PointsCampaign
	err := r.db.SelectContext(ctx, &campaigns, `
		SELECT id, name, description, campaign_type, points_reward, max_awards, start_date, end_date, active, created_at, updated_at
		FROM points_campaigns
		WHERE active = true 
		AND start_date <= NOW()
		AND (end_date IS NULL OR end_date > NOW())
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get active campaigns: %w", err)
	}
	return campaigns, nil
}

func (r *PointsRepository) CreateCampaignAward(ctx context.Context, tx *sqlx.Tx, award *models.CampaignAward) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO campaign_awards (id, campaign_id, user_id, transaction_id, awarded_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (campaign_id, user_id) DO NOTHING
	`, award.ID, award.CampaignID, award.UserID, award.TransactionID, award.AwardedAt)
	if err != nil {
		return fmt.Errorf("failed to create campaign award: %w", err)
	}
	return nil
}

func (r *PointsRepository) GetCampaignAwardCount(ctx context.Context, campaignID string) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		FROM campaign_awards
		WHERE campaign_id = $1
	`, campaignID)
	if err != nil {
		return 0, fmt.Errorf("failed to get campaign award count: %w", err)
	}
	return count, nil
}
