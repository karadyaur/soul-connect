package services

import (
	"context"
	"errors"

	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/repository"
	"soul-connect/sc-post/internal/utils"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(ctx context.Context, input models.AddCommentInput) (*models.Comment, error) {
	if input.Content == "" {
		return nil, errors.New("content is required")
	}
	postID, err := utils.UUIDFromString(input.PostID)
	if err != nil {
		return nil, err
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.CreateComment(ctx, db.CreateCommentParams{
		PostID:  postID,
		UserID:  userID,
		Content: input.Content,
	})
	if err != nil {
		return nil, err
	}

	model := commentFromDB(created)
	return &model, nil
}

func (s *CommentService) ListCommentsByPost(ctx context.Context, postID string) ([]models.Comment, error) {
	parsed, err := utils.UUIDFromString(postID)
	if err != nil {
		return nil, err
	}
	comments, err := s.repo.GetCommentsByPostID(ctx, parsed)
	if err != nil {
		return nil, err
	}
	return commentsFromDB(comments), nil
}
