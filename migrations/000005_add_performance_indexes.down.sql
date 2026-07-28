-- Rollback performance indexes

DROP INDEX IF EXISTS idx_posts_user_created;
DROP INDEX IF EXISTS idx_posts_visibility_created;
DROP INDEX IF EXISTS idx_messages_conversation_created;
DROP INDEX IF EXISTS idx_profiles_class_year;
DROP INDEX IF EXISTS idx_connections_status_created;
DROP INDEX IF EXISTS idx_connections_user_status;
DROP INDEX IF EXISTS idx_notifications_user_read;
DROP INDEX IF EXISTS idx_reactions_post_user;
DROP INDEX IF EXISTS idx_comments_post_created;
DROP INDEX IF EXISTS idx_business_location_status;
DROP INDEX IF EXISTS idx_jobs_deadline_status;
DROP INDEX IF EXISTS idx_merchant_offers_validity_status;
DROP INDEX IF EXISTS idx_sponsorships_active;
DROP INDEX IF EXISTS idx_group_members_group_role;
