-- Points and Rewards System for Partner Integrations (e.g., Naivas)
-- This migration adds tables for wallet management, transactions, redemptions, and partners

BEGIN;

-- Partners table - for merchant integrations (Naivas, etc.) - CREATE FIRST
CREATE TYPE integration_mode AS ENUM ('api', 'voucher');

CREATE TABLE partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    integration_mode integration_mode NOT NULL DEFAULT 'voucher',
    api_base_url TEXT,
    api_key_encrypted TEXT, -- Store encrypted in production
    callback_secret TEXT, -- For webhook signature verification
    logo_url TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    points_to_currency_rate DECIMAL(10, 2) DEFAULT 0.10, -- e.g., 100 points = KES 10
    min_points_to_redeem INTEGER DEFAULT 100,
    max_daily_redemptions INTEGER DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Points wallet table - one per user
CREATE TABLE points_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0,
    held_balance BIGINT NOT NULL DEFAULT 0, -- Points held during pending redemptions
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Points transaction table - records all point movements
CREATE TYPE transaction_type AS ENUM ('earn', 'redeem', 'hold', 'release', 'adjust', 'reversal');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'cancelled');

CREATE TABLE points_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES points_wallets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    amount BIGINT NOT NULL, -- Positive for earn, negative for spend
    balance_after BIGINT NOT NULL, -- Wallet balance after this transaction
    reference TEXT, -- Campaign ID, order ID, voucher ID, etc.
    partner_id UUID REFERENCES partners(id),
    status transaction_status NOT NULL DEFAULT 'pending',
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

-- Points redemption table - tracks redemption attempts
CREATE TYPE redemption_status AS ENUM ('requested', 'reserved', 'sent_to_partner', 'confirmed', 'failed', 'cancelled');

CREATE TABLE points_redemptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES points_transactions(id) ON DELETE CASCADE,
    partner_id UUID REFERENCES partners(id),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    redeem_method integration_mode NOT NULL,
    amount_points BIGINT NOT NULL,
    currency_value DECIMAL(10, 2),
    partner_reference TEXT, -- External reference from partner system
    voucher_code TEXT,
    voucher_expiry TIMESTAMPTZ,
    status redemption_status NOT NULL DEFAULT 'requested',
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Partner reconciliation table - for daily reconciliation reports
CREATE TABLE partner_reconciliation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID REFERENCES partners(id),
    partner_report_id TEXT,
    report_date DATE NOT NULL,
    total_redemptions INTEGER DEFAULT 0,
    total_points_redeemed BIGINT DEFAULT 0,
    total_currency_value DECIMAL(15, 2) DEFAULT 0,
    discrepancies JSONB,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Referral system table
CREATE TABLE referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_id UUID REFERENCES users(id) ON DELETE CASCADE, -- NULL until signup
    referral_code TEXT UNIQUE NOT NULL,
    points_awarded INTEGER DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, completed, expired
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Gamification: User badges and achievements
CREATE TYPE badge_type AS ENUM ('verified_alumni', 'super_referrer', 'early_adopter', 'top_contributor', 'mentor', 'donor');

CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_type badge_type NOT NULL,
    earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,
    UNIQUE(user_id, badge_type)
);

-- Gamification: User tiers/levels
CREATE TYPE user_tier AS ENUM ('bronze', 'silver', 'gold', 'platinum', 'diamond');

CREATE TABLE user_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier user_tier NOT NULL DEFAULT 'bronze',
    points_earned_lifetime BIGINT DEFAULT 0,
    current_level_progress INTEGER DEFAULT 0,
    promoted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Points campaigns (for admins to create point-earning opportunities)
CREATE TABLE points_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    campaign_type TEXT NOT NULL, -- signup_bonus, referral, event_attendance, social_share, etc.
    points_reward INTEGER NOT NULL,
    max_awards INTEGER, -- NULL for unlimited
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Campaign awards tracking
CREATE TABLE campaign_awards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES points_campaigns(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID REFERENCES points_transactions(id),
    awarded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(campaign_id, user_id)
);

-- Indexes for performance
CREATE INDEX idx_points_wallets_user_id ON points_wallets(user_id);
CREATE INDEX idx_points_wallets_balance ON points_wallets(balance);

CREATE INDEX idx_points_tx_user_id ON points_transactions(user_id);
CREATE INDEX idx_points_tx_wallet_id ON points_transactions(wallet_id);
CREATE INDEX idx_points_tx_status ON points_transactions(status);
CREATE INDEX idx_points_tx_type ON points_transactions(type);
CREATE INDEX idx_points_tx_created ON points_transactions(created_at DESC);
CREATE INDEX idx_points_tx_reference ON points_transactions(reference);

CREATE INDEX idx_partners_active ON partners(active);
CREATE INDEX idx_partners_mode ON partners(integration_mode);

CREATE INDEX idx_redemptions_user_id ON points_redemptions(user_id);
CREATE INDEX idx_redemptions_partner_id ON points_redemptions(partner_id);
CREATE INDEX idx_redemptions_status ON points_redemptions(status);
CREATE INDEX idx_redemptions_created ON points_redemptions(created_at DESC);
CREATE INDEX idx_redemptions_voucher ON points_redemptions(voucher_code);

CREATE INDEX idx_reconciliation_partner ON partner_reconciliation(partner_id);
CREATE INDEX idx_reconciliation_date ON partner_reconciliation(report_date);

CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);
CREATE INDEX idx_referrals_referred ON referrals(referred_id);
CREATE INDEX idx_referrals_code ON referrals(referral_code);
CREATE INDEX idx_referrals_status ON referrals(status);

CREATE INDEX idx_user_badges_user ON user_badges(user_id);
CREATE INDEX idx_user_badges_type ON user_badges(badge_type);

CREATE INDEX idx_user_tiers_user ON user_tiers(user_id);
CREATE INDEX idx_user_tiers_tier ON user_tiers(tier);

CREATE INDEX idx_campaigns_active ON points_campaigns(active);
CREATE INDEX idx_campaigns_dates ON points_campaigns(start_date, end_date);

CREATE INDEX idx_campaign_awards_campaign ON campaign_awards(campaign_id);
CREATE INDEX idx_campaign_awards_user ON campaign_awards(user_id);

-- Triggers for updated_at
CREATE TRIGGER update_points_wallets_updated_at BEFORE UPDATE ON points_wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_partners_updated_at BEFORE UPDATE ON partners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_redemptions_updated_at BEFORE UPDATE ON points_redemptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_tiers_updated_at BEFORE UPDATE ON user_tiers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_campaigns_updated_at BEFORE UPDATE ON points_campaigns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
