-- Rollback points and rewards system

BEGIN;

DROP TRIGGER IF EXISTS update_campaigns_updated_at ON points_campaigns;
DROP TRIGGER IF EXISTS update_user_tiers_updated_at ON user_tiers;
DROP TRIGGER IF EXISTS update_redemptions_updated_at ON points_redemptions;
DROP TRIGGER IF EXISTS update_partners_updated_at ON partners;
DROP TRIGGER IF EXISTS update_points_wallets_updated_at ON points_wallets;

DROP INDEX IF EXISTS idx_campaign_awards_user;
DROP INDEX IF EXISTS idx_campaign_awards_campaign;
DROP INDEX IF EXISTS idx_campaigns_dates;
DROP INDEX IF EXISTS idx_campaigns_active;
DROP INDEX IF EXISTS idx_user_tiers_tier;
DROP INDEX IF EXISTS idx_user_tiers_user;
DROP INDEX IF EXISTS idx_user_badges_type;
DROP INDEX IF EXISTS idx_user_badges_user;
DROP INDEX IF EXISTS idx_referrals_status;
DROP INDEX IF EXISTS idx_referrals_code;
DROP INDEX IF EXISTS idx_referrals_referred;
DROP INDEX IF EXISTS idx_referrals_referrer;
DROP INDEX IF EXISTS idx_reconciliation_date;
DROP INDEX IF EXISTS idx_reconciliation_partner;
DROP INDEX IF EXISTS idx_redemptions_voucher;
DROP INDEX IF EXISTS idx_redemptions_created;
DROP INDEX IF EXISTS idx_redemptions_status;
DROP INDEX IF EXISTS idx_redemptions_partner_id;
DROP INDEX IF EXISTS idx_redemptions_user_id;
DROP INDEX IF EXISTS idx_partners_mode;
DROP INDEX IF EXISTS idx_partners_active;
DROP INDEX IF EXISTS idx_points_tx_reference;
DROP INDEX IF EXISTS idx_points_tx_created;
DROP INDEX IF EXISTS idx_points_tx_type;
DROP INDEX IF EXISTS idx_points_tx_status;
DROP INDEX IF EXISTS idx_points_tx_wallet_id;
DROP INDEX IF EXISTS idx_points_tx_user_id;
DROP INDEX IF EXISTS idx_points_wallets_balance;
DROP INDEX IF EXISTS idx_points_wallets_user_id;

DROP TABLE IF EXISTS campaign_awards;
DROP TABLE IF EXISTS points_campaigns;
DROP TABLE IF EXISTS user_tiers;
DROP TABLE IF EXISTS user_badges;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS partner_reconciliation;
DROP TABLE IF EXISTS points_redemptions;
DROP TABLE IF EXISTS partners;
DROP TABLE IF EXISTS points_transactions;
DROP TABLE IF EXISTS points_wallets;

DROP TYPE IF EXISTS user_tier;
DROP TYPE IF EXISTS badge_type;
DROP TYPE IF EXISTS redemption_status;
DROP TYPE IF EXISTS integration_mode;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;

COMMIT;
