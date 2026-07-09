package services

import (
	"context"
	"fmt"
	"time"

	"starehian-society-platform/pkg/database"
	"starehian-society-platform/pkg/logger"
)

type AnalyticsService struct {
	db     *database.DB
	logger *logger.Logger
}

type AnalyticsMetrics struct {
	TotalUsers          int       `json:"total_users"`
	ActiveUsers         int       `json:"active_users"`
	TotalPosts          int       `json:"total_posts"`
	TotalConnections    int       `json:"total_connections"`
	TotalMessages       int       `json:"total_messages"`
	SignupsToday        int       `json:"signups_today"`
	SignupsThisWeek     int       `json:"signups_this_week"`
	SignupsThisMonth    int       `json:"signups_this_month"`
	PostsToday          int       `json:"posts_today"`
	MessagesToday       int       `json:"messages_today"`
	DAU                 int       `json:"dau"`
	MAU                 int       `json:"mau"`
}

type CohortAnalytics struct {
	ClassYear int `json:"class_year"`
	UserCount int `json:"user_count"`
	PostCount int `json:"post_count"`
}

type TimeSeriesData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func NewAnalyticsService(db *database.DB, logger *logger.Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger,
	}
}

// GetDashboardMetrics retrieves overall dashboard metrics
func (s *AnalyticsService) GetDashboardMetrics(ctx context.Context) (*AnalyticsMetrics, error) {
	metrics := &AnalyticsMetrics{}

	// Total users
	err := s.db.GetContext(ctx, &metrics.TotalUsers, `SELECT COUNT(*) FROM users WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Active users (last 7 days)
	err = s.db.GetContext(ctx, &metrics.ActiveUsers, `
		SELECT COUNT(DISTINCT user_id) 
		FROM (
			SELECT user_id FROM posts WHERE created_at >= NOW() - INTERVAL '7 days'
			UNION
			SELECT user_id FROM messages WHERE created_at >= NOW() - INTERVAL '7 days'
		) AS active_users
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// Total posts
	err = s.db.GetContext(ctx, &metrics.TotalPosts, `SELECT COUNT(*) FROM posts`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total posts: %w", err)
	}

	// Total connections
	err = s.db.GetContext(ctx, &metrics.TotalConnections, `SELECT COUNT(*) FROM connections WHERE status = 'accepted'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total connections: %w", err)
	}

	// Total messages
	err = s.db.GetContext(ctx, &metrics.TotalMessages, `SELECT COUNT(*) FROM messages`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total messages: %w", err)
	}

	// Signups today
	err = s.db.GetContext(ctx, &metrics.SignupsToday, `
		SELECT COUNT(*) FROM users 
		WHERE DATE(created_at) = CURRENT_DATE
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get signups today: %w", err)
	}

	// Signups this week
	err = s.db.GetContext(ctx, &metrics.SignupsThisWeek, `
		SELECT COUNT(*) FROM users 
		WHERE created_at >= DATE_TRUNC('week', CURRENT_DATE)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get signups this week: %w", err)
	}

	// Signups this month
	err = s.db.GetContext(ctx, &metrics.SignupsThisMonth, `
		SELECT COUNT(*) FROM users 
		WHERE DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get signups this month: %w", err)
	}

	// Posts today
	err = s.db.GetContext(ctx, &metrics.PostsToday, `
		SELECT COUNT(*) FROM posts 
		WHERE DATE(created_at) = CURRENT_DATE
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts today: %w", err)
	}

	// Messages today
	err = s.db.GetContext(ctx, &metrics.MessagesToday, `
		SELECT COUNT(*) FROM messages 
		WHERE DATE(created_at) = CURRENT_DATE
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages today: %w", err)
	}

	// DAU (Daily Active Users)
	err = s.db.GetContext(ctx, &metrics.DAU, `
		SELECT COUNT(DISTINCT user_id) 
		FROM (
			SELECT user_id FROM posts WHERE DATE(created_at) = CURRENT_DATE
			UNION
			SELECT user_id FROM messages WHERE DATE(created_at) = CURRENT_DATE
		) AS daily_active
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get DAU: %w", err)
	}

	// MAU (Monthly Active Users)
	err = s.db.GetContext(ctx, &metrics.MAU, `
		SELECT COUNT(DISTINCT user_id) 
		FROM (
			SELECT user_id FROM posts WHERE created_at >= NOW() - INTERVAL '30 days'
			UNION
			SELECT user_id FROM messages WHERE created_at >= NOW() - INTERVAL '30 days'
		) AS monthly_active
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get MAU: %w", err)
	}

	return metrics, nil
}

// GetCohortAnalytics retrieves analytics by class year
func (s *AnalyticsService) GetCohortAnalytics(ctx context.Context) ([]*CohortAnalytics, error) {
	var cohorts []*CohortAnalytics
	query := `
		SELECT 
			p.class_year,
			COUNT(DISTINCT p.user_id) as user_count,
			COUNT(DISTINCT ps.id) as post_count
		FROM profiles p
		LEFT JOIN posts ps ON p.user_id = ps.user_id
		WHERE p.class_year IS NOT NULL
		GROUP BY p.class_year
		ORDER BY p.class_year DESC
		LIMIT 20
	`

	err := s.db.SelectContext(ctx, &cohorts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cohort analytics: %w", err)
	}

	return cohorts, nil
}

// GetSignupsOverTime retrieves signup trends
func (s *AnalyticsService) GetSignupsOverTime(ctx context.Context, days int) ([]*TimeSeriesData, error) {
	var data []*TimeSeriesData
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users
		WHERE created_at >= NOW() - INTERVAL '1 day' * $1
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	err := s.db.SelectContext(ctx, &data, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get signups over time: %w", err)
	}

	return data, nil
}

// GetEngagementOverTime retrieves engagement trends
func (s *AnalyticsService) GetEngagementOverTime(ctx context.Context, days int) ([]*TimeSeriesData, error) {
	var data []*TimeSeriesData
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM (
			SELECT created_at FROM posts WHERE created_at >= NOW() - INTERVAL '1 day' * $1
			UNION ALL
			SELECT created_at FROM messages WHERE created_at >= NOW() - INTERVAL '1 day' * $1
		) AS activity
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	err := s.db.SelectContext(ctx, &data, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement over time: %w", err)
	}

	return data, nil
}

// GetTopContent retrieves most popular posts
func (s *AnalyticsService) GetTopContent(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	var content []map[string]interface{}
	query := `
		SELECT 
			p.id,
			p.user_id,
			p.content,
			p.created_at,
			COUNT(DISTINCT r.id) as reaction_count,
			COUNT(DISTINCT c.id) as comment_count
		FROM posts p
		LEFT JOIN reactions r ON p.id = r.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY p.id
		ORDER BY reaction_count DESC, comment_count DESC
		LIMIT $1
	`

	err := s.db.SelectContext(ctx, &content, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top content: %w", err)
	}

	return content, nil
}
