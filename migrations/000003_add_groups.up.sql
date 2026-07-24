-- Groups table for chapters, career spaces, cohorts, and custom groups
CREATE TYPE group_type AS ENUM ('chapter', 'career', 'cohort', 'custom');
CREATE TYPE join_policy AS ENUM ('open', 'approval_required', 'auto');

CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type group_type NOT NULL DEFAULT 'custom',
    join_policy join_policy NOT NULL DEFAULT 'open',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Group members table
CREATE TYPE member_role AS ENUM ('member', 'moderator', 'admin');

CREATE TABLE group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role member_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id)
);

-- Modify messages table to support both group chat and 1:1 DMs
ALTER TABLE messages 
ADD COLUMN group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
ADD COLUMN recipient_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Add indexes for groups
CREATE INDEX idx_groups_type ON groups(type);
CREATE INDEX idx_groups_created_at ON groups(created_at DESC);

CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_role ON group_members(role);

CREATE INDEX idx_messages_group_id ON messages(group_id);
CREATE INDEX idx_messages_recipient_id ON messages(recipient_id);

-- Add trigger for groups updated_at
CREATE TRIGGER update_groups_updated_at BEFORE UPDATE ON groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
