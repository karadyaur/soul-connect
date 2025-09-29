package services

import (
	"context"

	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/repository"
	"soul-connect/sc-post/internal/utils"
)

type LikeService struct {
	repo repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) LikePost(ctx context.Context, input models.LikeInput) (int32, error) {
	postID, err := utils.UUIDFromString(input.TargetID)
	if err != nil {
		return 0, err
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return 0, err
	}
	if err := s.repo.CreateLikeForPost(ctx, db.CreateLikeForPostParams{PostID: postID, UserID: userID}); err != nil {
		return 0, err
	}
	count, err := s.repo.GetLikesCountForPost(ctx, postID)
	if err != nil {
		return 0, err
	}
	if count.Valid {
		return count.Int32, nil
	}
	return 0, nil
}

func (s *LikeService) UnlikePost(ctx context.Context, input models.LikeInput) (int32, error) {
	postID, err := utils.UUIDFromString(input.TargetID)
	if err != nil {
		return 0, err
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return 0, err
	}
	if err := s.repo.DeleteLikeForPost(ctx, db.DeleteLikeForPostParams{PostID: postID, UserID: userID}); err != nil {
		return 0, err
	}
	count, err := s.repo.GetLikesCountForPost(ctx, postID)
	if err != nil {
		return 0, err
	}
	if count.Valid {
		return count.Int32, nil
	}
	return 0, nil
}

func (s *LikeService) LikeComment(ctx context.Context, input models.LikeInput) (int32, error) {
	commentID, err := utils.UUIDFromString(input.TargetID)
	if err != nil {
		return 0, err
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return 0, err
	}
	if err := s.repo.CreateLikeForComment(ctx, db.CreateLikeForCommentParams{CommentID: commentID, UserID: userID}); err != nil {
		return 0, err
	}
	count, err := s.repo.GetLikesCountForComment(ctx, commentID)
	if err != nil {
		return 0, err
	}
	if count.Valid {
		return count.Int32, nil
	}
	return 0, nil
}

func (s *LikeService) UnlikeComment(ctx context.Context, input models.LikeInput) (int32, error) {
	commentID, err := utils.UUIDFromString(input.TargetID)
	if err != nil {
		return 0, err
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return 0, err
	}
	if err := s.repo.DeleteLikeForComment(ctx, db.DeleteLikeForCommentParams{CommentID: commentID, UserID: userID}); err != nil {
		return 0, err
	}
	count, err := s.repo.GetLikesCountForComment(ctx, commentID)
	if err != nil {
		return 0, err
	}
	if count.Valid {
		return count.Int32, nil
	}
	return 0, nil
}
