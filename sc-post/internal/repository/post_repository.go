package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "soul-connect/sc-post/internal/db/sqlc"
)

type PostRepository interface {
	CreatePost(ctx context.Context, arg db.CreatePostParams) (db.Post, error)
	GetPostByID(ctx context.Context, id pgtype.UUID) (db.Post, error)
	GetPostsWithCommentsAndLikes(ctx context.Context) ([]db.GetPostsWithCommentsAndLikesRow, error)
	GetPostsByLabel(ctx context.Context, labelID pgtype.UUID) ([]db.Post, error)
	UpdatePost(ctx context.Context, arg db.UpdatePostParams) error
	DeletePost(ctx context.Context, id pgtype.UUID) error
}

type postRepository struct {
	queries db.Querier
}

func NewPostRepository(queries db.Querier) PostRepository {
	return &postRepository{queries: queries}
}

func (r *postRepository) CreatePost(ctx context.Context, arg db.CreatePostParams) (db.Post, error) {
	return r.queries.CreatePost(ctx, arg)
}

func (r *postRepository) GetPostByID(ctx context.Context, id pgtype.UUID) (db.Post, error) {
	return r.queries.GetPostByID(ctx, id)
}

func (r *postRepository) GetPostsWithCommentsAndLikes(ctx context.Context) ([]db.GetPostsWithCommentsAndLikesRow, error) {
	return r.queries.GetPostsWithCommentsAndLikes(ctx)
}

func (r *postRepository) GetPostsByLabel(ctx context.Context, labelID pgtype.UUID) ([]db.Post, error) {
	return r.queries.GetPostsByLabel(ctx, labelID)
}

func (r *postRepository) UpdatePost(ctx context.Context, arg db.UpdatePostParams) error {
	return r.queries.UpdatePost(ctx, arg)
}

func (r *postRepository) DeletePost(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeletePost(ctx, id)
}
