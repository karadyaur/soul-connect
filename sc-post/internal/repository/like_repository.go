package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "soul-connect/sc-post/internal/db/sqlc"
)

type LikeRepository interface {
	CreateLikeForPost(ctx context.Context, arg db.CreateLikeForPostParams) error
	DeleteLikeForPost(ctx context.Context, arg db.DeleteLikeForPostParams) error
	CreateLikeForComment(ctx context.Context, arg db.CreateLikeForCommentParams) error
	DeleteLikeForComment(ctx context.Context, arg db.DeleteLikeForCommentParams) error
	GetLikesCountForPost(ctx context.Context, postID pgtype.UUID) (pgtype.Int4, error)
	GetLikesCountForComment(ctx context.Context, commentID pgtype.UUID) (pgtype.Int4, error)
}

type likeRepository struct {
	queries db.Querier
}

func NewLikeRepository(queries db.Querier) LikeRepository {
	return &likeRepository{queries: queries}
}

func (r *likeRepository) CreateLikeForPost(ctx context.Context, arg db.CreateLikeForPostParams) error {
	return r.queries.CreateLikeForPost(ctx, arg)
}

func (r *likeRepository) DeleteLikeForPost(ctx context.Context, arg db.DeleteLikeForPostParams) error {
	return r.queries.DeleteLikeForPost(ctx, arg)
}

func (r *likeRepository) CreateLikeForComment(ctx context.Context, arg db.CreateLikeForCommentParams) error {
	return r.queries.CreateLikeForComment(ctx, arg)
}

func (r *likeRepository) DeleteLikeForComment(ctx context.Context, arg db.DeleteLikeForCommentParams) error {
	return r.queries.DeleteLikeForComment(ctx, arg)
}

func (r *likeRepository) GetLikesCountForPost(ctx context.Context, postID pgtype.UUID) (pgtype.Int4, error) {
	return r.queries.GetLikesCountForPost(ctx, postID)
}

func (r *likeRepository) GetLikesCountForComment(ctx context.Context, commentID pgtype.UUID) (pgtype.Int4, error) {
	return r.queries.GetLikesCountForComment(ctx, commentID)
}
