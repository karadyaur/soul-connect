package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "soul-connect/sc-post/internal/db/sqlc"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, arg db.CreateCommentParams) (db.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID pgtype.UUID) ([]db.Comment, error)
}

type commentRepository struct {
	queries db.Querier
}

func NewCommentRepository(queries db.Querier) CommentRepository {
	return &commentRepository{queries: queries}
}

func (r *commentRepository) CreateComment(ctx context.Context, arg db.CreateCommentParams) (db.Comment, error) {
	return r.queries.CreateComment(ctx, arg)
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID pgtype.UUID) ([]db.Comment, error) {
	return r.queries.GetCommentsByPostID(ctx, postID)
}
