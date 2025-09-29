package services

import (
	"context"

	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/repository"
	"soul-connect/sc-post/internal/utils"
)

type LabelService struct {
	repo repository.LabelRepository
}

func NewLabelService(repo repository.LabelRepository) *LabelService {
	return &LabelService{repo: repo}
}

func (s *LabelService) ListLabels(ctx context.Context) ([]models.Label, error) {
	labels, err := s.repo.GetAllLabels(ctx)
	if err != nil {
		return nil, err
	}
	return labelsFromDB(labels), nil
}

func (s *LabelService) GetLabelsForPost(ctx context.Context, postID string) ([]models.Label, error) {
	parsed, err := utils.UUIDFromString(postID)
	if err != nil {
		return nil, err
	}
	labels, err := s.repo.GetLabelsForPost(ctx, parsed)
	if err != nil {
		return nil, err
	}
	return labelsFromDB(labels), nil
}

func (s *LabelService) AddLabelToPost(ctx context.Context, input models.LabelAssignmentInput) error {
	postID, err := utils.UUIDFromString(input.PostID)
	if err != nil {
		return err
	}
	labelID, err := utils.UUIDFromString(input.LabelID)
	if err != nil {
		return err
	}
	return s.repo.AddLabelToPost(ctx, db.AddLabelToPostParams{PostID: postID, LabelID: labelID})
}

func (s *LabelService) RemoveLabelFromPost(ctx context.Context, input models.LabelAssignmentInput) error {
	postID, err := utils.UUIDFromString(input.PostID)
	if err != nil {
		return err
	}
	labelID, err := utils.UUIDFromString(input.LabelID)
	if err != nil {
		return err
	}
	return s.repo.RemoveLabelFromPost(ctx, db.RemoveLabelFromPostParams{PostID: postID, LabelID: labelID})
}
