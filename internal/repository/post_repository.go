package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/pkg/database"
)

type PostRepository struct {
	db *database.DB
}

func NewPostRepository(db *database.DB) *PostRepository {
	return &PostRepository{db: db}
}

// CreatePost creates a new post
func (r *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	// Convert media URLs slice to JSONB
	mediaURLsJSON, err := json.Marshal(post.MediaURLs)
	if err != nil {
		return fmt.Errorf("failed to marshal media URLs: %w", err)
	}

	query := `
		INSERT INTO posts (id, user_id, content, media_urls, visibility)
		VALUES (:id, :user_id, :content, :media_urls, :visibility)
		RETURNING created_at, updated_at
	`
	
	args := map[string]interface{}{
		"id":         post.ID,
		"user_id":    post.UserID,
		"content":    post.Content,
		"media_urls": mediaURLsJSON,
		"visibility": post.Visibility,
	}
	
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&post.CreatedAt, &post.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan post: %w", err)
		}
	}
	
	return nil
}

// GetPost retrieves a post by ID
func (r *PostRepository) GetPost(ctx context.Context, postID string) (*models.Post, error) {
	var post models.Post
	var mediaURLsJSON []byte
	
	query := `
		SELECT id, user_id, content, media_urls, visibility, created_at, updated_at
		FROM posts
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &post, query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	
	// Unmarshal media URLs
	if mediaURLsJSON != nil {
		if err := json.Unmarshal(mediaURLsJSON, &post.MediaURLs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal media URLs: %w", err)
		}
	}
	
	return &post, nil
}

// UpdatePost updates a post
func (r *PostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	mediaURLsJSON, err := json.Marshal(post.MediaURLs)
	if err != nil {
		return fmt.Errorf("failed to marshal media URLs: %w", err)
	}

	query := `
		UPDATE posts 
		SET content = :content, media_urls = :media_urls, visibility = :visibility
		WHERE id = :id
		RETURNING updated_at
	`
	
	args := map[string]interface{}{
		"id":         post.ID,
		"content":    post.Content,
		"media_urls": mediaURLsJSON,
		"visibility": post.Visibility,
	}
	
	result, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("post not found")
	}
	
	return nil
}

// DeletePost deletes a post
func (r *PostRepository) DeletePost(ctx context.Context, postID string) error {
	query := `DELETE FROM posts WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("post not found")
	}
	
	return nil
}

// GetPostsByUser retrieves posts by a user using offset pagination
func (r *PostRepository) GetPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	query := `
		SELECT id, user_id, content, media_urls, visibility, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &posts, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by user: %w", err)
	}
	
	// Unmarshal media URLs for each post
	for _, post := range posts {
		if post.Content != nil {
			// Need to re-query to get media_urls as JSONB
			var mediaURLsJSON []byte
			mediaQuery := `SELECT media_urls FROM posts WHERE id = $1`
			err := r.db.GetContext(ctx, &mediaURLsJSON, mediaQuery, post.ID)
			if err == nil && mediaURLsJSON != nil {
				json.Unmarshal(mediaURLsJSON, &post.MediaURLs)
			}
		}
	}
	
	return posts, nil
}

// GetPostsByUserCursor retrieves posts by a user using cursor-based pagination
func (r *PostRepository) GetPostsByUserCursor(ctx context.Context, userID string, limit int, cursor *string) ([]*models.Post, *string, error) {
	var posts []*models.Post
	
	query := `
		SELECT id, user_id, content, media_urls, visibility, created_at, updated_at
		FROM posts
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	argCount := 1
	
	if cursor != nil {
		query += fmt.Sprintf(" AND created_at < (SELECT created_at FROM posts WHERE id = $%d)", argCount+1)
		args = append(args, *cursor)
		argCount++
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit+1) // Fetch one extra to check if there's more
		argCount++
	}
	
	err := r.db.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get posts by user cursor: %w", err)
	}
	
	var nextCursor *string
	if len(posts) > limit {
		posts = posts[:limit]
		nextCursor = &posts[len(posts)-1].ID
	}
	
	// Unmarshal media URLs for each post
	for _, post := range posts {
		if post.Content != nil {
			var mediaURLsJSON []byte
			mediaQuery := `SELECT media_urls FROM posts WHERE id = $1`
			err := r.db.GetContext(ctx, &mediaURLsJSON, mediaQuery, post.ID)
			if err == nil && mediaURLsJSON != nil {
				json.Unmarshal(mediaURLsJSON, &post.MediaURLs)
			}
		}
	}
	
	return posts, nextCursor, nil
}

// GetFeed retrieves posts for a user's feed using offset pagination
func (r *PostRepository) GetFeed(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	// Get posts from user and their connections
	query := `
		SELECT DISTINCT p.id, p.user_id, p.content, p.media_urls, p.visibility, p.created_at, p.updated_at
		FROM posts p
		WHERE p.user_id = $1
		   OR p.visibility = 'public'
		   OR (p.visibility = 'connections' AND EXISTS (
		       SELECT 1 FROM connections c 
		       WHERE (c.user_id = $1 AND c.connected_user_id = p.user_id AND c.status = 'accepted')
		       OR (c.connected_user_id = $1 AND c.user_id = p.user_id AND c.status = 'accepted')
		   ))
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &posts, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	
	return posts, nil
}

// GetFeedCursor retrieves posts for a user's feed using cursor-based pagination
func (r *PostRepository) GetFeedCursor(ctx context.Context, userID string, limit int, cursor *string) ([]*models.Post, *string, error) {
	var posts []*models.Post
	
	query := `
		SELECT DISTINCT p.id, p.user_id, p.content, p.media_urls, p.visibility, p.created_at, p.updated_at
		FROM posts p
		WHERE p.user_id = $1
		   OR p.visibility = 'public'
		   OR (p.visibility = 'connections' AND EXISTS (
		       SELECT 1 FROM connections c 
		       WHERE (c.user_id = $1 AND c.connected_user_id = p.user_id AND c.status = 'accepted')
		       OR (c.connected_user_id = $1 AND c.user_id = p.user_id AND c.status = 'accepted')
		   ))
	`
	
	args := []interface{}{userID}
	argCount := 1
	
	if cursor != nil {
		query += fmt.Sprintf(" AND p.created_at < (SELECT created_at FROM posts WHERE id = $%d)", argCount+1)
		args = append(args, *cursor)
		argCount++
	}
	
	query += " ORDER BY p.created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit+1) // Fetch one extra to check if there's more
		argCount++
	}
	
	err := r.db.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get feed cursor: %w", err)
	}
	
	var nextCursor *string
	if len(posts) > limit {
		posts = posts[:limit]
		nextCursor = &posts[len(posts)-1].ID
	}
	
	return posts, nextCursor, nil
}

// Comments
func (r *PostRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, user_id, content)
		VALUES (:id, :post_id, :user_id, :content)
		RETURNING created_at, updated_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, comment)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan comment: %w", err)
		}
	}
	
	return nil
}

func (r *PostRepository) GetComments(ctx context.Context, postID string, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := `
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &comments, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	
	return comments, nil
}

func (r *PostRepository) DeleteComment(ctx context.Context, commentID string) error {
	query := `DELETE FROM comments WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("comment not found")
	}
	
	return nil
}

// Reactions
func (r *PostRepository) CreateReaction(ctx context.Context, reaction *models.Reaction) error {
	query := `
		INSERT INTO reactions (id, post_id, user_id, type)
		VALUES (:id, :post_id, :user_id, :type)
		ON CONFLICT (post_id, user_id) DO UPDATE SET type = :type
		RETURNING created_at
	`
	
	rows, err := r.db.NamedQueryContext(ctx, query, reaction)
	if err != nil {
		return fmt.Errorf("failed to create reaction: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&reaction.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan reaction: %w", err)
		}
	}
	
	return nil
}

func (r *PostRepository) GetReactions(ctx context.Context, postID string) ([]*models.Reaction, error) {
	var reactions []*models.Reaction
	query := `
		SELECT id, post_id, user_id, type, created_at
		FROM reactions
		WHERE post_id = $1
		ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &reactions, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}
	
	return reactions, nil
}

func (r *PostRepository) DeleteReaction(ctx context.Context, postID, userID string) error {
	query := `DELETE FROM reactions WHERE post_id = $1 AND user_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("reaction not found")
	}
	
	return nil
}

func (r *PostRepository) GetUserReaction(ctx context.Context, postID, userID string) (*models.Reaction, error) {
	var reaction models.Reaction
	query := `
		SELECT id, post_id, user_id, type, created_at
		FROM reactions
		WHERE post_id = $1 AND user_id = $2
	`
	
	err := r.db.GetContext(ctx, &reaction, query, postID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user reaction: %w", err)
	}
	
	return &reaction, nil
}
