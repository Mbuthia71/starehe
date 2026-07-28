package models

import "time"

// PointsWallet represents a user's points wallet
type PointsWallet struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Balance     int64     `json:"balance" db:"balance"`
	HeldBalance int64     `json:"held_balance" db:"held_balance"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionType represents the type of points transaction
type TransactionType string

const (
	TransactionTypeEarn     TransactionType = "earn"
	TransactionTypeRedeem   TransactionType = "redeem"
	TransactionTypeHold     TransactionType = "hold"
	TransactionTypeRelease  TransactionType = "release"
	TransactionTypeAdjust   TransactionType = "adjust"
	TransactionTypeReversal TransactionType = "reversal"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// PointsTransaction represents a points transaction
type PointsTransaction struct {
	ID           string             `json:"id" db:"id"`
	WalletID     string             `json:"wallet_id" db:"wallet_id"`
	UserID       string             `json:"user_id" db:"user_id"`
	Type         TransactionType    `json:"type" db:"type"`
	Amount       int64              `json:"amount" db:"amount"`
	BalanceAfter int64              `json:"balance_after" db:"balance_after"`
	Reference    *string            `json:"reference,omitempty" db:"reference"`
	PartnerID    *string            `json:"partner_id,omitempty" db:"partner_id"`
	Status       TransactionStatus  `json:"status" db:"status"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	ProcessedAt  *time.Time         `json:"processed_at,omitempty" db:"processed_at"`
}

// IntegrationMode represents how a partner integrates
type IntegrationMode string

const (
	IntegrationModeAPI     IntegrationMode = "api"
	IntegrationModeVoucher IntegrationMode = "voucher"
)

// Partner represents a merchant partner (e.g., Naivas)
type Partner struct {
	ID                    string          `json:"id" db:"id"`
	Name                  string          `json:"name" db:"name"`
	Description           *string         `json:"description,omitempty" db:"description"`
	IntegrationMode       IntegrationMode `json:"integration_mode" db:"integration_mode"`
	APIBaseURL            *string         `json:"api_base_url,omitempty" db:"api_base_url"`
	APIKeyEncrypted       *string         `json:"-" db:"api_key_encrypted"` // Never expose in JSON
	CallbackSecret        *string         `json:"-" db:"callback_secret"`    // Never expose in JSON
	LogoURL               *string         `json:"logo_url,omitempty" db:"logo_url"`
	Active                bool            `json:"active" db:"active"`
	PointsToCurrencyRate  float64         `json:"points_to_currency_rate" db:"points_to_currency_rate"`
	MinPointsToRedeem     int             `json:"min_points_to_redeem" db:"min_points_to_redeem"`
	MaxDailyRedemptions   int             `json:"max_daily_redemptions" db:"max_daily_redemptions"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}

// RedemptionStatus represents the status of a redemption
type RedemptionStatus string

const (
	RedemptionStatusRequested     RedemptionStatus = "requested"
	RedemptionStatusReserved      RedemptionStatus = "reserved"
	RedemptionStatusSentToPartner RedemptionStatus = "sent_to_partner"
	RedemptionStatusConfirmed     RedemptionStatus = "confirmed"
	RedemptionStatusFailed        RedemptionStatus = "failed"
	RedemptionStatusCancelled     RedemptionStatus = "cancelled"
)

// PointsRedemption represents a redemption attempt
type PointsRedemption struct {
	ID               string          `json:"id" db:"id"`
	TransactionID    string          `json:"transaction_id" db:"transaction_id"`
	PartnerID        *string         `json:"partner_id,omitempty" db:"partner_id"`
	UserID           string          `json:"user_id" db:"user_id"`
	RedeemMethod     IntegrationMode `json:"redeem_method" db:"redeem_method"`
	AmountPoints     int64           `json:"amount_points" db:"amount_points"`
	CurrencyValue    *float64        `json:"currency_value,omitempty" db:"currency_value"`
	PartnerReference *string         `json:"partner_reference,omitempty" db:"partner_reference"`
	VoucherCode      *string         `json:"voucher_code,omitempty" db:"voucher_code"`
	VoucherExpiry    *time.Time      `json:"voucher_expiry,omitempty" db:"voucher_expiry"`
	Status           RedemptionStatus `json:"status" db:"status"`
	FailureReason    *string         `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// PartnerReconciliation represents a daily reconciliation report
type PartnerReconciliation struct {
	ID                   string                  `json:"id" db:"id"`
	PartnerID            *string                 `json:"partner_id,omitempty" db:"partner_id"`
	PartnerReportID      *string                 `json:"partner_report_id,omitempty" db:"partner_report_id"`
	ReportDate           time.Time               `json:"report_date" db:"report_date"`
	TotalRedemptions     int                     `json:"total_redemptions" db:"total_redemptions"`
	TotalPointsRedeemed  int64                   `json:"total_points_redeemed" db:"total_points_redeemed"`
	TotalCurrencyValue   float64                 `json:"total_currency_value" db:"total_currency_value"`
	Discrepancies        map[string]interface{}  `json:"discrepancies,omitempty" db:"discrepancies"`
	ProcessedAt          *time.Time              `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt            time.Time               `json:"created_at" db:"created_at"`
}

// Referral represents a referral record
type Referral struct {
	ID            string     `json:"id" db:"id"`
	ReferrerID    string     `json:"referrer_id" db:"referrer_id"`
	ReferredID    *string    `json:"referred_id,omitempty" db:"referred_id"`
	ReferralCode  string     `json:"referral_code" db:"referral_code"`
	PointsAwarded int        `json:"points_awarded" db:"points_awarded"`
	Status        string     `json:"status" db:"status"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// BadgeType represents the type of badge
type BadgeType string

const (
	BadgeTypeVerifiedAlumni BadgeType = "verified_alumni"
	BadgeTypeSuperReferrer  BadgeType = "super_referrer"
	BadgeTypeEarlyAdopter   BadgeType = "early_adopter"
	BadgeTypeTopContributor BadgeType = "top_contributor"
	BadgeTypeMentor         BadgeType = "mentor"
	BadgeTypeDonor          BadgeType = "donor"
)

// UserBadge represents a user's badge
type UserBadge struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	BadgeType BadgeType              `json:"badge_type" db:"badge_type"`
	EarnedAt  time.Time              `json:"earned_at" db:"earned_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

// UserTier represents a user's tier level
type UserTier string

const (
	UserTierBronze   UserTier = "bronze"
	UserTierSilver   UserTier = "silver"
	UserTierGold     UserTier = "gold"
	UserTierPlatinum UserTier = "platinum"
	UserTierDiamond  UserTier = "diamond"
)

// UserTierRecord represents a user's tier record
type UserTierRecord struct {
	ID                 string    `json:"id" db:"id"`
	UserID             string    `json:"user_id" db:"user_id"`
	Tier               UserTier  `json:"tier" db:"tier"`
	PointsEarnedLifetime int64    `json:"points_earned_lifetime" db:"points_earned_lifetime"`
	CurrentLevelProgress int      `json:"current_level_progress" db:"current_level_progress"`
	PromotedAt         *time.Time `json:"promoted_at,omitempty" db:"promoted_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// PointsCampaign represents a points campaign
type PointsCampaign struct {
	ID            string     `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	CampaignType  string     `json:"campaign_type" db:"campaign_type"`
	PointsReward  int        `json:"points_reward" db:"points_reward"`
	MaxAwards     *int       `json:"max_awards,omitempty" db:"max_awards"`
	StartDate     time.Time  `json:"start_date" db:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty" db:"end_date"`
	Active        bool       `json:"active" db:"active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CampaignAward represents a campaign award
type CampaignAward struct {
	ID            string     `json:"id" db:"id"`
	CampaignID    string     `json:"campaign_id" db:"campaign_id"`
	UserID        string     `json:"user_id" db:"user_id"`
	TransactionID *string    `json:"transaction_id,omitempty" db:"transaction_id"`
	AwardedAt     time.Time  `json:"awarded_at" db:"awarded_at"`
}

// Request/Response DTOs

// RedeemPointsRequest represents a redemption request
type RedeemPointsRequest struct {
	AmountPoints int64          `json:"amount_points" validate:"required,min=1"`
	PartnerID     string         `json:"partner_id" validate:"required"`
	Method        IntegrationMode `json:"method" validate:"required,oneof=api voucher"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// RedeemPointsResponse represents a redemption response
type RedeemPointsResponse struct {
	RedemptionID   string          `json:"redemption_id"`
	TransactionID  string          `json:"transaction_id"`
	Status         RedemptionStatus `json:"status"`
	AmountPoints   int64           `json:"amount_points"`
	CurrencyValue  float64         `json:"currency_value"`
	VoucherCode    *string         `json:"voucher_code,omitempty"`
	VoucherExpiry  *time.Time      `json:"voucher_expiry,omitempty"`
	Next           map[string]interface{} `json:"next,omitempty"`
}

// PointsBalanceResponse represents a points balance response
type PointsBalanceResponse struct {
	Balance            int64   `json:"balance"`
	HeldBalance        int64   `json:"held_balance"`
	Currency            string  `json:"currency"`
	PointsToCurrencyRate float64 `json:"points_to_currency_rate"`
	CurrencyValue       float64 `json:"currency_value"`
}

// PointsHistoryResponse represents a points history response
type PointsHistoryResponse struct {
	Transactions []PointsTransaction `json:"transactions"`
	Cursor       *string             `json:"cursor,omitempty"`
	HasMore      bool                `json:"has_more"`
}

// AdjustPointsRequest represents an admin adjustment request
type AdjustPointsRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required"`
	Reason  string `json:"reason" validate:"required"`
}

// CreatePartnerRequest represents a partner creation request
type CreatePartnerRequest struct {
	Name                 string          `json:"name" validate:"required"`
	Description          *string         `json:"description,omitempty"`
	IntegrationMode      IntegrationMode `json:"integration_mode" validate:"required,oneof=api voucher"`
	APIBaseURL           *string         `json:"api_base_url,omitempty"`
	APIKey               *string         `json:"api_key,omitempty"`
	PointsToCurrencyRate float64         `json:"points_to_currency_rate" validate:"required,min=0"`
	MinPointsToRedeem    int             `json:"min_points_to_redeem" validate:"required,min=1"`
	MaxDailyRedemptions  int             `json:"max_daily_redemptions" validate:"required,min=1"`
}

// CreateReferralRequest represents a referral creation request
type CreateReferralRequest struct {
	ReferrerID string `json:"referrer_id" validate:"required"`
}

// ReferralResponse represents a referral response
type ReferralResponse struct {
	ReferralCode string `json:"referral_code"`
	ReferralURL  string `json:"referral_url"`
	PointsReward int    `json:"points_reward"`
}
