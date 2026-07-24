package repository

import (
	"context"
	"database/sql"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type GroupRepository struct {
	db *database.DB
}

func NewGroupRepository(db *database.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// CreateGroup creates a new group
func (r *GroupRepository) CreateGroup(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (id, name, type, join_policy, description)
		VALUES (:id, :name, :type, :join_policy, :description)
		RETURNING created_at, updated_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, group)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&group.CreatedAt, &group.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan group: %w", err)
		}
	}
	
	return nil
}

// GetGroup retrieves a group by ID
func (r *GroupRepository) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, name, type, join_policy, description, created_at, updated_at
		FROM groups
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &group, query, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	
	return &group, nil
}

// GetGroupsByType retrieves groups by type
func (r *GroupRepository) GetGroupsByType(ctx context.Context, groupType models.GroupType, limit, offset int) ([]*models.Group, error) {
	var groups []*models.Group
	query := `
		SELECT id, name, type, join_policy, description, created_at, updated_at
		FROM groups
		WHERE type = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &groups, query, groupType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups by type: %w", err)
	}
	
	return groups, nil
}

// GetGroupsWithMemberCount retrieves groups with member count for a user
func (r *GroupRepository) GetGroupsWithMemberCount(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithMemberCount, error) {
	var groups []*models.GroupWithMemberCount
	query := `
		SELECT 
			g.id, g.name, g.type, g.join_policy, g.description, g.created_at, g.updated_at,
			COUNT(gm.user_id) as member_count,
			COALESCE(bool_or(gm.user_id = $1), false) as is_member,
			CASE WHEN bool_or(gm.user_id = $1) THEN 
				(SELECT role FROM group_members WHERE group_id = g.id AND user_id = $1 LIMIT 1)
			ELSE NULL END as user_role
		FROM groups g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		GROUP BY g.id, g.name, g.type, g.join_policy, g.description, g.created_at, g.updated_at
		ORDER BY g.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &groups, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups with member count: %w", err)
	}
	
	return groups, nil
}

// GetUserGroups retrieves groups that a user is a member of
func (r *GroupRepository) GetUserGroups(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithMemberCount, error) {
	var groups []*models.GroupWithMemberCount
	query := `
		SELECT 
			g.id, g.name, g.type, g.join_policy, g.description, g.created_at, g.updated_at,
			COUNT(gm.user_id) as member_count,
			true as is_member,
			mgm.role as user_role
		FROM groups g
		INNER JOIN group_members mgm ON g.id = mgm.group_id AND mgm.user_id = $1
		LEFT JOIN group_members gm ON g.id = gm.group_id
		GROUP BY g.id, g.name, g.type, g.join_policy, g.description, g.created_at, g.updated_at, mgm.role
		ORDER BY g.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &groups, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	
	return groups, nil
}

// AddGroupMember adds a member to a group
func (r *GroupRepository) AddGroupMember(ctx context.Context, member *models.GroupMember) error {
	query := `
		INSERT INTO group_members (id, group_id, user_id, role)
		VALUES (:id, :group_id, :user_id, :role)
		RETURNING joined_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, member)
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&member.JoinedAt); err != nil {
			return fmt.Errorf("failed to scan group member: %w", err)
		}
	}
	
	return nil
}

// GetGroupMembers retrieves members of a group
func (r *GroupRepository) GetGroupMembers(ctx context.Context, groupID string, limit, offset int) ([]*models.GroupMember, error) {
	var members []*models.GroupMember
	query := `
		SELECT id, group_id, user_id, role, joined_at
		FROM group_members
		WHERE group_id = $1
		ORDER BY joined_at ASC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &members, query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	
	return members, nil
}

// GetGroupMember retrieves a specific group member
func (r *GroupRepository) GetGroupMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	var member models.GroupMember
	query := `
		SELECT id, group_id, user_id, role, joined_at
		FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`
	
	err := r.db.GetContext(ctx, &member, query, groupID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group member: %w", err)
	}
	
	return &member, nil
}

// UpdateGroupMemberRole updates a member's role in a group
func (r *GroupRepository) UpdateGroupMemberRole(ctx context.Context, groupID, userID string, role models.MemberRole) error {
	query := `
		UPDATE group_members
		SET role = $1
		WHERE group_id = $2 AND user_id = $3
	`
	
	result, err := r.db.ExecContext(ctx, query, role, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to update group member role: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("group member not found")
	}
	
	return nil
}

// RemoveGroupMember removes a member from a group
func (r *GroupRepository) RemoveGroupMember(ctx context.Context, groupID, userID string) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove group member: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("group member not found")
	}
	
	return nil
}

// GetCohortGroupByYear retrieves a cohort group for a specific graduation year
func (r *GroupRepository) GetCohortGroupByYear(ctx context.Context, classYear int) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, name, type, join_policy, description, created_at, updated_at
		FROM groups
		WHERE type = 'cohort' AND name = $1
		LIMIT 1
	`
	
	groupName := fmt.Sprintf("Class of %d", classYear)
	err := r.db.GetContext(ctx, &group, query, groupName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cohort group: %w", err)
	}
	
	return &group, nil
}

// GetCareerGroupByField retrieves a career group for a specific field
func (r *GroupRepository) GetCareerGroupByField(ctx context.Context, careerField string) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, name, type, join_policy, description, created_at, updated_at
		FROM groups
		WHERE type = 'career' AND name = $1
		LIMIT 1
	`
	
	err := r.db.GetContext(ctx, &group, query, careerField)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get career group: %w", err)
	}
	
	return &group, nil
}
