package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
	userRepo  *repository.UserRepository
	logger    *logger.Logger
}

type CreateGroupRequest struct {
	Name        string           `json:"name"`
	Type        models.GroupType `json:"type"`
	JoinPolicy  models.JoinPolicy `json:"join_policy"`
	Description string           `json:"description,omitempty"`
}

type UpdateGroupRequest struct {
	Name        *string          `json:"name,omitempty"`
	JoinPolicy  *models.JoinPolicy `json:"join_policy,omitempty"`
	Description *string          `json:"description,omitempty"`
}

func NewGroupService(
	groupRepo *repository.GroupRepository,
	userRepo *repository.UserRepository,
	logger *logger.Logger,
) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(ctx context.Context, creatorID string, req *CreateGroupRequest) (*models.Group, error) {
	group := &models.Group{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Type:        req.Type,
		JoinPolicy:  req.JoinPolicy,
		Description: &req.Description,
	}

	err := s.groupRepo.CreateGroup(ctx, group)
	if err != nil {
		s.logger.Errorf("Failed to create group: %v", err)
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator as admin
	member := &models.GroupMember{
		ID:      uuid.New().String(),
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    models.MemberRoleAdmin,
	}

	err = s.groupRepo.AddGroupMember(ctx, member)
	if err != nil {
		s.logger.Errorf("Failed to add creator as admin: %v", err)
		return nil, fmt.Errorf("failed to add creator as admin: %w", err)
	}

	s.logger.Infof("Group created: %s by %s", group.ID, creatorID)
	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	group, err := s.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		s.logger.Errorf("Failed to get group: %v", err)
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	return group, nil
}

// ListGroups lists all groups with member counts
func (s *GroupService) ListGroups(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithMemberCount, error) {
	groups, err := s.groupRepo.GetGroupsWithMemberCount(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to list groups: %v", err)
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	return groups, nil
}

// ListMyGroups lists groups that the user is a member of
func (s *GroupService) ListMyGroups(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithMemberCount, error) {
	groups, err := s.groupRepo.GetUserGroups(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to list user groups: %v", err)
		return nil, fmt.Errorf("failed to list user groups: %w", err)
	}

	return groups, nil
}

// JoinGroup allows a user to join a group
func (s *GroupService) JoinGroup(ctx context.Context, userID, groupID string) error {
	group, err := s.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check if already a member
	existingMember, err := s.groupRepo.GetGroupMember(ctx, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}

	if existingMember != nil {
		return fmt.Errorf("already a member of this group")
	}

	// Check join policy
	if group.JoinPolicy == models.JoinPolicyApprovalRequired {
		return fmt.Errorf("approval required - please contact group admin")
	}

	// Add member
	member := &models.GroupMember{
		ID:      uuid.New().String(),
		GroupID: groupID,
		UserID:  userID,
		Role:    models.MemberRoleMember,
	}

	err = s.groupRepo.AddGroupMember(ctx, member)
	if err != nil {
		s.logger.Errorf("Failed to add member to group: %v", err)
		return fmt.Errorf("failed to join group: %w", err)
	}

	s.logger.Infof("User %s joined group %s", userID, groupID)
	return nil
}

// LeaveGroup allows a user to leave a group
func (s *GroupService) LeaveGroup(ctx context.Context, userID, groupID string) error {
	err := s.groupRepo.RemoveGroupMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Errorf("Failed to remove member from group: %v", err)
		return fmt.Errorf("failed to leave group: %w", err)
	}

	s.logger.Infof("User %s left group %s", userID, groupID)
	return nil
}

// GetGroupMembers retrieves members of a group
func (s *GroupService) GetGroupMembers(ctx context.Context, groupID string, limit, offset int) ([]*models.GroupMember, error) {
	members, err := s.groupRepo.GetGroupMembers(ctx, groupID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get group members: %v", err)
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	return members, nil
}

// AddMember adds a member to a group (admin only)
func (s *GroupService) AddMember(ctx context.Context, adminID, groupID, newMemberID string) error {
	// Check if requester is admin
	member, err := s.groupRepo.GetGroupMember(ctx, groupID, adminID)
	if err != nil {
		return fmt.Errorf("failed to check admin status: %w", err)
	}

	if member == nil || member.Role != models.MemberRoleAdmin {
		return fmt.Errorf("only admins can add members")
	}

	// Check if already a member
	existingMember, err := s.groupRepo.GetGroupMember(ctx, groupID, newMemberID)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}

	if existingMember != nil {
		return fmt.Errorf("user is already a member")
	}

	// Add new member
	newMember := &models.GroupMember{
		ID:      uuid.New().String(),
		GroupID: groupID,
		UserID:  newMemberID,
		Role:    models.MemberRoleMember,
	}

	err = s.groupRepo.AddGroupMember(ctx, newMember)
	if err != nil {
		s.logger.Errorf("Failed to add member: %v", err)
		return fmt.Errorf("failed to add member: %w", err)
	}

	s.logger.Infof("Member %s added to group %s by %s", newMemberID, groupID, adminID)
	return nil
}

// RemoveMember removes a member from a group (admin only)
func (s *GroupService) RemoveMember(ctx context.Context, adminID, groupID, memberID string) error {
	// Check if requester is admin
	member, err := s.groupRepo.GetGroupMember(ctx, groupID, adminID)
	if err != nil {
		return fmt.Errorf("failed to check admin status: %w", err)
	}

	if member == nil || member.Role != models.MemberRoleAdmin {
		return fmt.Errorf("only admins can remove members")
	}

	err = s.groupRepo.RemoveGroupMember(ctx, groupID, memberID)
	if err != nil {
		s.logger.Errorf("Failed to remove member: %v", err)
		return fmt.Errorf("failed to remove member: %w", err)
	}

	s.logger.Infof("Member %s removed from group %s by %s", memberID, groupID, adminID)
	return nil
}

// UpdateMemberRole updates a member's role (admin only)
func (s *GroupService) UpdateMemberRole(ctx context.Context, adminID, groupID, memberID string, newRole models.MemberRole) error {
	// Check if requester is admin
	member, err := s.groupRepo.GetGroupMember(ctx, groupID, adminID)
	if err != nil {
		return fmt.Errorf("failed to check admin status: %w", err)
	}

	if member == nil || member.Role != models.MemberRoleAdmin {
		return fmt.Errorf("only admins can update member roles")
	}

	err = s.groupRepo.UpdateGroupMemberRole(ctx, groupID, memberID, newRole)
	if err != nil {
		s.logger.Errorf("Failed to update member role: %v", err)
		return fmt.Errorf("failed to update member role: %w", err)
	}

	s.logger.Infof("Member %s role updated to %s in group %s by %s", memberID, newRole, groupID, adminID)
	return nil
}

// GetOrCreateCohortGroup gets or creates a cohort group for a graduation year
func (s *GroupService) GetOrCreateCohortGroup(ctx context.Context, classYear int) (*models.Group, error) {
	// Try to get existing cohort group
	group, err := s.groupRepo.GetCohortGroupByYear(ctx, classYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get cohort group: %w", err)
	}

	if group != nil {
		return group, nil
	}

	// Create new cohort group
	groupName := fmt.Sprintf("Class of %d", classYear)
	group = &models.Group{
		ID:          uuid.New().String(),
		Name:        groupName,
		Type:        models.GroupTypeCohort,
		JoinPolicy:  models.JoinPolicyAuto,
		Description: &groupName,
	}

	err = s.groupRepo.CreateGroup(ctx, group)
	if err != nil {
		s.logger.Errorf("Failed to create cohort group: %v", err)
		return nil, fmt.Errorf("failed to create cohort group: %w", err)
	}

	s.logger.Infof("Cohort group created: %s for year %d", group.ID, classYear)
	return group, nil
}

// GetOrCreateCareerGroup gets or creates a career group for a field
func (s *GroupService) GetOrCreateCareerGroup(ctx context.Context, careerField string) (*models.Group, error) {
	// Try to get existing career group
	group, err := s.groupRepo.GetCareerGroupByField(ctx, careerField)
	if err != nil {
		return nil, fmt.Errorf("failed to get career group: %w", err)
	}

	if group != nil {
		return group, nil
	}

	// Create new career group
	group = &models.Group{
		ID:          uuid.New().String(),
		Name:        careerField,
		Type:        models.GroupTypeCareer,
		JoinPolicy:  models.JoinPolicyOpen,
		Description: &careerField,
	}

	err = s.groupRepo.CreateGroup(ctx, group)
	if err != nil {
		s.logger.Errorf("Failed to create career group: %w", err)
		return nil, fmt.Errorf("failed to create career group: %w", err)
	}

	s.logger.Infof("Career group created: %s for field %s", group.ID, careerField)
	return group, nil
}
