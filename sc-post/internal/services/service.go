package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/events"
	"soul-connect/sc-post/internal/repository"
)

type Services struct {
	Posts    *PostService
	Comments *CommentService
	Likes    *LikeService
	Labels   *LabelService
}

func NewServices(pool *pgxpool.Pool, publisher events.PostEventPublisher) *Services {
	queries := db.New(pool)
	repo := repository.New(queries)

	return &Services{
		Posts:    NewPostService(repo.Posts, repo.Labels, repo.Comments, publisher),
		Comments: NewCommentService(repo.Comments),
		Likes:    NewLikeService(repo.Likes),
		Labels:   NewLabelService(repo.Labels),
	}
}
