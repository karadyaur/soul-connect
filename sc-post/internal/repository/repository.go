package repository

import db "soul-connect/sc-post/internal/db/sqlc"

type Repository struct {
	Posts    PostRepository
	Comments CommentRepository
	Likes    LikeRepository
	Labels   LabelRepository
}

func New(queries db.Querier) *Repository {
	return &Repository{
		Posts:    NewPostRepository(queries),
		Comments: NewCommentRepository(queries),
		Likes:    NewLikeRepository(queries),
		Labels:   NewLabelRepository(queries),
	}
}
