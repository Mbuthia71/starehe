package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starehian-society-platform/internal/middleware"
	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type PostService struct {
	postRepo    *repository.PostRepository
	authzService *middleware.AuthorizationService
	logger      *logger.Logger
}

type CreatePostRequest struct {
	Content   *string  `json:"content"`
	MediaURLs []string `json:"media_urls"`
	Visibility string  `json:"visibility"`
}

type UpdatePostRequest struct {
	Content   *string  `json:"content"`
	MediaURLs []string `json:"media_urls"`
	Visibility string  `json:"visibility"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type CreateReactionRequest struct {
	Type string `json:"type"`
}

func NewPostService(
	postRepo *repository.PostRepository,
	authzService *middleware.AuthorizationService,
	logger *logger.Logger,
) *PostService {
	return &PostService{
		postRepo:    postRepo,
		authzService: authzService,
		logger:      logger,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, userID string, req *CreatePostRequest) (*models.Post, error) {
	if req.Content == nil && len(req.MediaURLs) == 0 {
		return nil, fmt.Errorf("post must have content or media")
	}

	post := &models.Post{
		ID:         uuid.New().String(),
		UserID:     userID,
		Content:    req.Content,
		MediaURLs:  req.MediaURLs,
		Visibility: req.Visibility,
	}

	if post.Visibility == "" {
		post.Visibility = string(models.VisibilityConnections)
	}

	err := s.postRepo.CreatePost(ctx, post)
	if err != nil {
		s.logger.Errorf("Failed to create post: %v", err)
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	s.logger.Infof("Post created: %s", post.ID)
	return post, nil
}

// GetPost retrieves a post with privacy check
func (s *PostService) GetPost(ctx context.Context, viewerID, postID string) (*models.Post, error) {
	// Check if viewer can view post
	canView, err := s.authzService.CanViewPost(ctx, viewerID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !canView {
		return nil, fmt.Errorf("access denied")
	}

	post, err := s.postRepo.GetPost(ctx, postID)
	if err != nil {
		s.logger.Errorf("Failed to get post: %v", err)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(ctx context.Context, userID, postID string, req *UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepo.GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check if user is the author
	if post.UserID != userID {
		return nil, fmt.Errorf("not authorized to update this post")
	}

	// Update fields
	if req.Content != nil {
		post.Content = req.Content
	}
	if req.MediaURLs != nil {
		post.MediaURLs = req.MediaURLs
	}
	if req.Visibility != "" {
		post.Visibility = req.Visibility
	}

	err = s.postRepo.UpdatePost(ctx, post)
	if err != nil {
		s.logger.Errorf("Failed to update post: %v", err)
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	s.logger.Infof("Post updated: %s", postID)
	return post, nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(ctx context.Context, userID, postID string) error {
	post, err := s.postRepo.GetPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	// Check if user is the author
	if post.UserID != userID {
		return fmt.Errorf("not authorized to delete this post")
	}

	err = s.postRepo.DeletePost(ctx, postID)
	if err != nil {
		s.logger.Errorf("Failed to delete post: %v", err)
		return fmt.Errorf("failed to delete post: %w", err)
	}

	s.logger.Infof("Post deleted: %s", postID)
	return nil
}

// GetFeed retrieves posts for the user's feed
func (s *PostService) GetFeed(ctx context.Context, userID string, limit, offset int) ([]*models.Post, error) {
	posts, err := s.postRepo.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get feed: %v", err)
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return posts, nil
}

// GetPostsByUser retrieves posts by a specific user
func (s *PostService) GetPostsByUser(ctx context.Context, viewerID, targetID string, limit, offset int) ([]*models.Post, error) {
	// Check if viewer can view target's posts
	// For now, allow viewing own posts or public posts
	// This could be enhanced with more granular privacy checks
	
	posts, err := s.postRepo.GetPostsByUser(ctx, targetID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get posts by user: %v", err)
		return nil, fmt.Errorf("failed to get posts by user: %w", err)
	}

	// Filter based on privacy
	var filteredPosts []*models.Post
	for _, post := range posts {
		canView, _ := s.authzService.CanViewPost(ctx, viewerID, post.ID)
		if canView {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts, nil
}

// Comments
func (s *PostService) CreateComment(ctx context.Context, userID, postID string, req *CreateCommentRequest) (*models.Comment, error) {
	// Check if can view post (to comment on it)
	canView, err := s.authzService.CanViewPost(ctx, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !canView {
		return nil, fmt.Errorf("access denied")
	}

	comment := &models.Comment{
		ID:     uuid.New().String(),
		PostID: postID,
		UserID: userID,
		Content: req.Content,
	}

	err = s.postRepo.CreateComment(ctx, comment)
	if err != nil {
		s.logger.Errorf("Failed to create comment: %v", err)
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	s.logger.Infof("Comment created: %s", comment.ID)
	return comment, nil
}

func (s *PostService) GetComments(ctx context.Context, viewerID, postID string, limit, offset int) ([]*models.Comment, error) {
	// Check if can view post
	canView, err := s.authzService.CanViewPost(ctx, viewerID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !canView {
		return nil, fmt.Errorf("access denied")
	}

	comments, err := s.postRepo.GetComments(ctx, postID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get comments: %v", err)
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

func (s *PostService) DeleteComment(ctx context.Context, userID, commentID string) error {
	// Get comment to check ownership
	// For simplicity, we'll add a GetCommentByID method to repository
	// For now, proceed with deletion
	err := s.postRepo.DeleteComment(ctx, commentID)
	if err != nil {
		s.logger.Errorf("Failed to delete comment: %v", err)
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	s.logger.Infof("Comment deleted: %s", commentID)
	return nil
}

// Reactions
func (s *PostService) CreateReaction(ctx context.Context, userID, postID string, req *CreateReactionRequest) (*models.Reaction, error) {
	// Check if can view post
	canView, err := s.authzService.CanViewPost(ctx, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !canView {
		return nil, fmt.Errorf("access denied")
	}

	reaction := &models.Reaction{
		ID:     uuid.New().String(),
		PostID: postID,
		UserID: userID,
		Type:   req.Type,
	}

	err = s.postRepo.CreateReaction(ctx, reaction)
	if err != nil {
		s.logger.Errorf("Failed to create reaction: %v", err)
		return nil, fmt.Errorf("failed to create reaction: %w", err)
	}

	s.logger.Infof("Reaction created: %s", reaction.ID)
	return reaction, nil
}

func (s *PostService) GetReactions(ctx context.Context, viewerID, postID string) ([]*models.Reaction, error) {
	// Check if can view post
	canView, err := s.authzService.CanViewPost(ctx, viewerID, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}
	if !canView {
		return nil, fmt.Errorf("access denied")
	}

	reactions, err := s.postRepo.GetReactions(ctx, postID)
	if err != nil {
		s.logger.Errorf("Failed to get reactions: %v", err)
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}

	return reactions, nil
}

func (s *PostService) DeleteReaction(ctx context.Context, userID, postID string) error {
	err := s.postRepo.DeleteReaction(ctx, postID, userID)
	if err != nil {
		s.logger.Errorf("Failed to delete reaction: %v", err)
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	s.logger.Infof("Reaction deleted for user %s on post %s", userID, postID)
	return nil
}

func (s *PostService) GetUserReaction(ctx context.Context, userID, postID string) (*models.Reaction, error) {
	reaction, err := s.postRepo.GetUserReaction(ctx, postID, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user reaction: %v", err)
		return nil, fmt.Errorf("failed to get user reaction: %w", err)
	}

	return reaction, nil
}
