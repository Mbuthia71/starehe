-- Drop trigger for groups
DROP TRIGGER IF EXISTS update_groups_updated_at ON groups;

-- Drop indexes
DROP INDEX IF EXISTS idx_messages_recipient_id;
DROP INDEX IF EXISTS idx_messages_group_id;
DROP INDEX IF EXISTS idx_group_members_role;
DROP INDEX IF EXISTS idx_group_members_user_id;
DROP INDEX IF EXISTS idx_group_members_group_id;
DROP INDEX IF EXISTS idx_groups_created_at;
DROP INDEX IF EXISTS idx_groups_type;

-- Modify messages table to remove group_id and recipient_id
ALTER TABLE messages 
DROP COLUMN IF EXISTS recipient_id,
DROP COLUMN IF EXISTS group_id;

-- Drop group members table
DROP TABLE IF EXISTS group_members;

-- Drop groups table
DROP TABLE IF EXISTS groups;

-- Drop custom types
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS join_policy;
DROP TYPE IF EXISTS group_type;
