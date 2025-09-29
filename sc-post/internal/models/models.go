package models

import "time"

type Label struct {
	ID   string
	Name string
}

type Comment struct {
	ID         string
	PostID     string
	UserID     string
	Content    string
	LikesCount int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Post struct {
	ID          string
	UserID      string
	Title       string
	Description string
	LikesCount  int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Labels      []Label
}

type PostSummary struct {
	Post          Post
	TotalComments int64
	TotalLikes    int64
}

type CreatePostInput struct {
	UserID      string
	Title       string
	Description string
	LabelIDs    []string
}

type UpdatePostInput struct {
	ID          string
	Title       *string
	Description *string
}

type AddCommentInput struct {
	PostID  string
	UserID  string
	Content string
}

type LikeInput struct {
	TargetID string
	UserID   string
}

type LabelAssignmentInput struct {
	PostID  string
	LabelID string
}
