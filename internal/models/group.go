package models

import "time"

type GroupType string
type JoinPolicy string
type MemberRole string

const (
	GroupTypeChapter GroupType = "chapter"
	GroupTypeCareer  GroupType = "career"
	GroupTypeCohort  GroupType = "cohort"
	GroupTypeCustom  GroupType = "custom"
)

const (
	JoinPolicyOpen           JoinPolicy = "open"
	JoinPolicyApprovalRequired JoinPolicy = "approval_required"
	JoinPolicyAuto           JoinPolicy = "auto"
)

const (
	MemberRoleMember    MemberRole = "member"
	MemberRoleModerator MemberRole = "moderator"
	MemberRoleAdmin     MemberRole = "admin"
)

type Group struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Type        GroupType  `json:"type" db:"type"`
	JoinPolicy  JoinPolicy `json:"join_policy" db:"join_policy"`
	Description *string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type GroupMember struct {
	ID        string     `json:"id" db:"id"`
	GroupID   string     `json:"group_id" db:"group_id"`
	UserID    string     `json:"user_id" db:"user_id"`
	Role      MemberRole `json:"role" db:"role"`
	JoinedAt  time.Time  `json:"joined_at" db:"joined_at"`
}

type GroupWithMemberCount struct {
	ID           string     `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Type         GroupType  `json:"type" db:"type"`
	JoinPolicy   JoinPolicy `json:"join_policy" db:"join_policy"`
	Description  *string    `json:"description,omitempty" db:"description"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	MemberCount  int        `json:"member_count" db:"member_count"`
	IsMember     bool       `json:"is_member" db:"is_member"`
	UserRole     *string    `json:"user_role,omitempty" db:"user_role"`
}
