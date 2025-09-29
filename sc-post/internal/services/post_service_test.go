package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/events"
	"soul-connect/sc-post/internal/models"
)

type stubPostRepo struct {
	CreatePostFn                   func(ctx context.Context, arg db.CreatePostParams) (db.Post, error)
	GetPostByIDFn                  func(ctx context.Context, id pgtype.UUID) (db.Post, error)
	GetPostsWithCommentsAndLikesFn func(ctx context.Context) ([]db.GetPostsWithCommentsAndLikesRow, error)
	GetPostsByLabelFn              func(ctx context.Context, labelID pgtype.UUID) ([]db.Post, error)
	UpdatePostFn                   func(ctx context.Context, arg db.UpdatePostParams) error
	DeletePostFn                   func(ctx context.Context, id pgtype.UUID) error
}

func (s *stubPostRepo) CreatePost(ctx context.Context, arg db.CreatePostParams) (db.Post, error) {
	return s.CreatePostFn(ctx, arg)
}

func (s *stubPostRepo) GetPostByID(ctx context.Context, id pgtype.UUID) (db.Post, error) {
	return s.GetPostByIDFn(ctx, id)
}

func (s *stubPostRepo) GetPostsWithCommentsAndLikes(ctx context.Context) ([]db.GetPostsWithCommentsAndLikesRow, error) {
	return s.GetPostsWithCommentsAndLikesFn(ctx)
}

func (s *stubPostRepo) GetPostsByLabel(ctx context.Context, labelID pgtype.UUID) ([]db.Post, error) {
	return s.GetPostsByLabelFn(ctx, labelID)
}

func (s *stubPostRepo) UpdatePost(ctx context.Context, arg db.UpdatePostParams) error {
	return s.UpdatePostFn(ctx, arg)
}

func (s *stubPostRepo) DeletePost(ctx context.Context, id pgtype.UUID) error {
	return s.DeletePostFn(ctx, id)
}

type stubLabelRepo struct {
	AddLabelToPostFn   func(ctx context.Context, arg db.AddLabelToPostParams) error
	GetLabelsForPostFn func(ctx context.Context, postID pgtype.UUID) ([]db.Label, error)
}

func (s *stubLabelRepo) AddLabelToPost(ctx context.Context, arg db.AddLabelToPostParams) error {
	return s.AddLabelToPostFn(ctx, arg)
}

func (s *stubLabelRepo) RemoveLabelFromPost(ctx context.Context, arg db.RemoveLabelFromPostParams) error {
	return nil
}

func (s *stubLabelRepo) GetLabelsForPost(ctx context.Context, postID pgtype.UUID) ([]db.Label, error) {
	return s.GetLabelsForPostFn(ctx, postID)
}

func (s *stubLabelRepo) GetAllLabels(ctx context.Context) ([]db.Label, error) {
	return nil, nil
}

type stubCommentRepo struct {
	GetCommentsByPostIDFn func(ctx context.Context, postID pgtype.UUID) ([]db.Comment, error)
}

func (s *stubCommentRepo) CreateComment(ctx context.Context, arg db.CreateCommentParams) (db.Comment, error) {
	return db.Comment{}, nil
}

func (s *stubCommentRepo) GetCommentsByPostID(ctx context.Context, postID pgtype.UUID) ([]db.Comment, error) {
	return s.GetCommentsByPostIDFn(ctx, postID)
}

type recorderPublisher struct {
	events.PostEventPublisher
	called bool
}

func (r *recorderPublisher) PublishPostCreated(ctx context.Context, post models.Post) error {
	r.called = true
	return nil
}

func TestPostService_CreatePostPublishesEvent(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	var addCalls []db.AddLabelToPostParams
	postRepo := &stubPostRepo{
		CreatePostFn: func(_ context.Context, arg db.CreatePostParams) (db.Post, error) {
			require.Equal(t, userID.String(), uuid.UUID(arg.UserID.Bytes[:]).String())
			return db.Post{
				ID:          toPgUUID(postID),
				UserID:      toPgUUID(userID),
				Title:       arg.Title,
				Description: arg.Description,
				LikesCount:  pgtype.Int4{Int32: 0, Valid: true},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			}, nil
		},
	}
	labelRepo := &stubLabelRepo{
		AddLabelToPostFn: func(_ context.Context, arg db.AddLabelToPostParams) error {
			addCalls = append(addCalls, arg)
			return nil
		},
		GetLabelsForPostFn: func(_ context.Context, _ pgtype.UUID) ([]db.Label, error) {
			return []db.Label{{ID: toPgUUID(uuid.New()), Name: "Happy"}}, nil
		},
	}
	commentRepo := &stubCommentRepo{}
	publisher := &recorderPublisher{}

	service := NewPostService(postRepo, labelRepo, commentRepo, publisher)

	created, err := service.CreatePost(ctx, models.CreatePostInput{
		UserID:      userID.String(),
		Title:       "Test title",
		Description: "Description",
		LabelIDs:    []string{uuid.New().String()},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.True(t, publisher.called)
	require.Len(t, addCalls, 1)
	require.Equal(t, "Test title", created.Title)
	require.Equal(t, userID.String(), created.UserID)
	require.Len(t, created.Labels, 1)
}

func TestPostService_ListPostsAggregatesLabels(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	postRepo := &stubPostRepo{
		GetPostsWithCommentsAndLikesFn: func(context.Context) ([]db.GetPostsWithCommentsAndLikesRow, error) {
			return []db.GetPostsWithCommentsAndLikesRow{{
				PostID:          toPgUUID(postID),
				PostUserID:      toPgUUID(userID),
				PostTitle:       "Title",
				PostDescription: pgtype.Text{String: "Desc", Valid: true},
				PostLikes:       pgtype.Int4{Int32: 2, Valid: true},
				TotalComments:   3,
				TotalLikes:      2,
			}}, nil
		},
	}
	labelRepo := &stubLabelRepo{
		GetLabelsForPostFn: func(context.Context, pgtype.UUID) ([]db.Label, error) {
			return []db.Label{{ID: toPgUUID(uuid.New()), Name: "Calm"}}, nil
		},
	}
	commentRepo := &stubCommentRepo{}
	publisher := &recorderPublisher{}

	service := NewPostService(postRepo, labelRepo, commentRepo, publisher)

	summaries, err := service.ListPosts(ctx, nil)
	require.NoError(t, err)
	require.Len(t, summaries, 1)
	summary := summaries[0]
	require.Equal(t, int64(3), summary.TotalComments)
	require.Equal(t, int64(2), summary.TotalLikes)
	require.Equal(t, "Title", summary.Post.Title)
	require.Equal(t, userID.String(), summary.Post.UserID)
	require.Len(t, summary.Post.Labels, 1)
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	var b [16]byte
	copy(b[:], id[:])
	return pgtype.UUID{Bytes: b, Valid: true}
}
