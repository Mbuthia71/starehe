-- Performance indexes for 10,000+ users scale
-- These indexes optimize the most common query patterns

-- Composite index for user feed queries (user_id + created_at)
-- This improves: "Get posts by user" and "Get feed" with pagination
CREATE INDEX IF NOT EXISTS idx_posts_user_created ON posts(user_id, created_at DESC);

-- Composite index for feed filtering by visibility
-- This improves: "Get public posts" and "Get connection posts"
CREATE INDEX IF NOT EXISTS idx_posts_visibility_created ON posts(visibility, created_at DESC);

-- Composite index for message pagination in conversations
-- This improves: "Get messages in conversation" with cursor pagination
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at DESC);

-- Index on profiles.class_year for cohort filtering
-- This improves: "Search by class year" and "Cohort analytics"
CREATE INDEX IF NOT EXISTS idx_profiles_class_year ON profiles(class_year);

-- Composite index for connection status queries
-- This improves: "Get pending connections" and "Get accepted connections"
CREATE INDEX IF NOT EXISTS idx_connections_status_created ON connections(status, created_at DESC);

-- Composite index for user + status in connections
-- This improves: "Get user's connections by status"
CREATE INDEX IF NOT EXISTS idx_connections_user_status ON connections(user_id, status);

-- Index for notification read status filtering
-- This improves: "Get unread notifications"
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, read_at);

-- Composite index for post reactions lookup
-- This improves: "Get user's reaction on post"
CREATE INDEX IF NOT EXISTS idx_reactions_post_user ON reactions(post_id, user_id);

-- Index for comment pagination on posts
-- This improves: "Get comments for post" with pagination
CREATE INDEX IF NOT EXISTS idx_comments_post_created ON comments(post_id, created_at DESC);

-- Index for business listings by location and status
-- This improves: "Search businesses by location"
CREATE INDEX IF NOT EXISTS idx_business_location_status ON business_listings(location, status);

-- Index for jobs by deadline and status
-- This improves: "Get active jobs" and "Get closing soon jobs"
CREATE INDEX IF NOT EXISTS idx_jobs_deadline_status ON jobs(application_deadline, status);

-- Index for merchant offers by validity period
-- This improves: "Get active offers" and "Get expiring offers"
CREATE INDEX IF NOT EXISTS idx_merchant_offers_validity_status ON merchant_offers(valid_from, valid_until, status);

-- Partial index for active sponsorships only
-- This improves: "Get active sponsorships" query
CREATE INDEX IF NOT EXISTS idx_sponsorships_active ON sponsorships(start_date, end_date) WHERE status = 'active';

-- Index for group members by role
-- This improves: "Get group admins" and "Get group moderators"
CREATE INDEX IF NOT EXISTS idx_group_members_group_role ON group_members(group_id, role);
