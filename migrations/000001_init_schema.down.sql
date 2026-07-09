-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_connections_updated_at ON connections;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_reports_updated_at ON reports;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_profiles_full_name_trgm;
DROP INDEX IF EXISTS idx_profiles_career_trgm;
DROP INDEX IF EXISTS idx_profiles_location_trgm;

DROP INDEX IF EXISTS idx_alumni_roster_phone;
DROP INDEX IF EXISTS idx_alumni_roster_house;
DROP INDEX IF EXISTS idx_alumni_roster_class_year;

DROP INDEX IF EXISTS idx_admin_actions_created_at;
DROP INDEX IF EXISTS idx_admin_actions_admin_id;

DROP INDEX IF EXISTS idx_reports_target;
DROP INDEX IF EXISTS idx_reports_status;
DROP INDEX IF EXISTS idx_reports_reporter_id;

DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_read_at;
DROP INDEX IF EXISTS idx_notifications_user_id;

DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_conversation_id;

DROP INDEX IF EXISTS idx_conversation_members_user_id;
DROP INDEX IF EXISTS idx_conversation_members_conversation_id;

DROP INDEX IF EXISTS idx_reactions_user_id;
DROP INDEX IF EXISTS idx_reactions_post_id;

DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;

DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_visibility;
DROP INDEX IF EXISTS idx_posts_user_id;

DROP INDEX IF EXISTS idx_blocks_blocked_id;
DROP INDEX IF EXISTS idx_blocks_blocker_id;

DROP INDEX IF EXISTS idx_connections_status;
DROP INDEX IF EXISTS idx_connections_connected_user_id;
DROP INDEX IF EXISTS idx_connections_user_id;

DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;

-- Drop tables
DROP TABLE IF EXISTS alumni_roster;
DROP TABLE IF EXISTS admin_actions;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversation_members;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS blocks;
DROP TABLE IF EXISTS connections;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS pg_trgm;
