package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/events"
	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/repository"
	"soul-connect/sc-post/internal/utils"
)

type PostService struct {
	postRepo    repository.PostRepository
	labelRepo   repository.LabelRepository
	commentRepo repository.CommentRepository
	publisher   events.PostEventPublisher
}

func NewPostService(postRepo repository.PostRepository, labelRepo repository.LabelRepository, commentRepo repository.CommentRepository, publisher events.PostEventPublisher) *PostService {
	return &PostService{postRepo: postRepo, labelRepo: labelRepo, commentRepo: commentRepo, publisher: publisher}
}

func (s *PostService) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	userID, err := utils.UUIDFromString(input.UserID)
	if err != nil {
		return nil, err
	}

	params := db.CreatePostParams{
		UserID:      userID,
		Title:       input.Title,
		Description: utils.TextFromString(input.Description),
	}

	created, err := s.postRepo.CreatePost(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, labelID := range input.LabelIDs {
		if err := s.attachLabel(ctx, created.ID, labelID); err != nil {
			return nil, err
		}
	}

	labels, err := s.labelRepo.GetLabelsForPost(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	post := postFromDB(created, labels)
	if s.publisher != nil {
		if err := s.publisher.PublishPostCreated(ctx, post); err != nil {
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostService) GetPost(ctx context.Context, id string) (*models.Post, error) {
	postID, err := utils.UUIDFromString(id)
	if err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	labels, err := s.labelRepo.GetLabelsForPost(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	result := postFromDB(post, labels)
	return &result, nil
}

func (s *PostService) ListPosts(ctx context.Context, labelIDs []string) ([]models.PostSummary, error) {
	if len(labelIDs) == 0 {
		rows, err := s.postRepo.GetPostsWithCommentsAndLikes(ctx)
		if err != nil {
			return nil, err
		}
		summaries := make([]models.PostSummary, 0, len(rows))
		for _, row := range rows {
			labels, err := s.labelRepo.GetLabelsForPost(ctx, row.PostID)
			if err != nil {
				return nil, err
			}
			summary := postSummaryFromRow(row)
			summary.Post.Labels = labelsFromDB(labels)
			summaries = append(summaries, summary)
		}
		return summaries, nil
	}

	summaries := make([]models.PostSummary, 0)
	seen := make(map[string]struct{})
	for _, labelIDStr := range labelIDs {
		labelID, err := utils.UUIDFromString(labelIDStr)
		if err != nil {
			return nil, err
		}
		posts, err := s.postRepo.GetPostsByLabel(ctx, labelID)
		if err != nil {
			return nil, err
		}
		for _, post := range posts {
			id := utils.UUIDToString(post.ID)
			if _, exists := seen[id]; exists {
				continue
			}
			labels, err := s.labelRepo.GetLabelsForPost(ctx, post.ID)
			if err != nil {
				return nil, err
			}
			comments, err := s.commentRepo.GetCommentsByPostID(ctx, post.ID)
			if err != nil {
				return nil, err
			}
			model := postFromDB(post, labels)
			summary := models.PostSummary{
				Post:          model,
				TotalComments: int64(len(comments)),
				TotalLikes:    int64(model.LikesCount),
			}
			summaries = append(summaries, summary)
			seen[id] = struct{}{}
		}
	}
	return summaries, nil
}

func (s *PostService) UpdatePost(ctx context.Context, input models.UpdatePostInput) (*models.Post, error) {
	postID, err := utils.UUIDFromString(input.ID)
	if err != nil {
		return nil, err
	}

	existing, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	updateParams := db.UpdatePostParams{
		Title:       existing.Title,
		Description: existing.Description,
		ID:          postID,
	}

	if input.Title != nil {
		updateParams.Title = *input.Title
	}

	if input.Description != nil {
		updateParams.Description = utils.NullableTextFromPointer(input.Description)
	}

	if err := s.postRepo.UpdatePost(ctx, updateParams); err != nil {
		return nil, err
	}

	updated, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	labels, err := s.labelRepo.GetLabelsForPost(ctx, updated.ID)
	if err != nil {
		return nil, err
	}

	model := postFromDB(updated, labels)
	return &model, nil
}

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	postID, err := utils.UUIDFromString(id)
	if err != nil {
		return err
	}
	return s.postRepo.DeletePost(ctx, postID)
}

func (s *PostService) attachLabel(ctx context.Context, postID pgtype.UUID, labelID string) error {
	if labelID == "" {
		return nil
	}
	parsed, err := utils.UUIDFromString(labelID)
	if err != nil {
		return err
	}
	return s.labelRepo.AddLabelToPost(ctx, db.AddLabelToPostParams{
		LabelID: parsed,
		PostID:  postID,
	})
}
