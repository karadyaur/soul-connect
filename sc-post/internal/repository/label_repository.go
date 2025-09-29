package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "soul-connect/sc-post/internal/db/sqlc"
)

type LabelRepository interface {
	AddLabelToPost(ctx context.Context, arg db.AddLabelToPostParams) error
	RemoveLabelFromPost(ctx context.Context, arg db.RemoveLabelFromPostParams) error
	GetLabelsForPost(ctx context.Context, postID pgtype.UUID) ([]db.Label, error)
	GetAllLabels(ctx context.Context) ([]db.Label, error)
}

type labelRepository struct {
	queries db.Querier
}

func NewLabelRepository(queries db.Querier) LabelRepository {
	return &labelRepository{queries: queries}
}

func (r *labelRepository) AddLabelToPost(ctx context.Context, arg db.AddLabelToPostParams) error {
	return r.queries.AddLabelToPost(ctx, arg)
}

func (r *labelRepository) RemoveLabelFromPost(ctx context.Context, arg db.RemoveLabelFromPostParams) error {
	return r.queries.RemoveLabelFromPost(ctx, arg)
}

func (r *labelRepository) GetLabelsForPost(ctx context.Context, postID pgtype.UUID) ([]db.Label, error) {
	return r.queries.GetLabelsForPost(ctx, postID)
}

func (r *labelRepository) GetAllLabels(ctx context.Context) ([]db.Label, error) {
	return r.queries.GetAllLabels(ctx)
}
