package server

import (
	"context"
	"time"

	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/services"
	postpb "soul-connect/sc-post/pkg/postpb"
)

type PostServer struct {
	postpb.UnimplementedPostServiceServer
	services *services.Services
}

func NewPostServer(services *services.Services) *PostServer {
	return &PostServer{services: services}
}

func (s *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.CreatePostResponse, error) {
	input := models.CreatePostInput{
		UserID:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		LabelIDs:    req.LabelIds,
	}
	post, err := s.services.Posts.CreatePost(ctx, input)
	if err != nil {
		return nil, err
	}
	return &postpb.CreatePostResponse{Post: toProtoPost(*post)}, nil
}

func (s *PostServer) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.GetPostResponse, error) {
	post, err := s.services.Posts.GetPost(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	comments, err := s.services.Comments.ListCommentsByPost(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &postpb.GetPostResponse{
		Post:     toProtoPost(*post),
		Comments: toProtoComments(comments),
	}, nil
}

func (s *PostServer) ListPosts(ctx context.Context, req *postpb.ListPostsRequest) (*postpb.ListPostsResponse, error) {
	posts, err := s.services.Posts.ListPosts(ctx, req.LabelIds)
	if err != nil {
		return nil, err
	}
	summaries := make([]*postpb.PostSummary, 0, len(posts))
	for _, summary := range posts {
		protoSummary := &postpb.PostSummary{
			Post:          toProtoPost(summary.Post),
			TotalComments: summary.TotalComments,
			TotalLikes:    summary.TotalLikes,
		}
		summaries = append(summaries, protoSummary)
	}
	return &postpb.ListPostsResponse{Posts: summaries}, nil
}

func (s *PostServer) AddComment(ctx context.Context, req *postpb.AddCommentRequest) (*postpb.AddCommentResponse, error) {
	comment, err := s.services.Comments.AddComment(ctx, models.AddCommentInput{
		PostID:  req.PostId,
		UserID:  req.UserId,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &postpb.AddCommentResponse{Comment: toProtoComment(*comment)}, nil
}

func (s *PostServer) ListComments(ctx context.Context, req *postpb.ListCommentsRequest) (*postpb.ListCommentsResponse, error) {
	comments, err := s.services.Comments.ListCommentsByPost(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	return &postpb.ListCommentsResponse{Comments: toProtoComments(comments)}, nil
}

func (s *PostServer) LikePost(ctx context.Context, req *postpb.LikePostRequest) (*postpb.LikeCountResponse, error) {
	likes, err := s.services.Likes.LikePost(ctx, models.LikeInput{TargetID: req.PostId, UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return &postpb.LikeCountResponse{LikesCount: likes}, nil
}

func (s *PostServer) UnlikePost(ctx context.Context, req *postpb.UnlikePostRequest) (*postpb.LikeCountResponse, error) {
	likes, err := s.services.Likes.UnlikePost(ctx, models.LikeInput{TargetID: req.PostId, UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return &postpb.LikeCountResponse{LikesCount: likes}, nil
}

func (s *PostServer) LikeComment(ctx context.Context, req *postpb.LikeCommentRequest) (*postpb.LikeCountResponse, error) {
	likes, err := s.services.Likes.LikeComment(ctx, models.LikeInput{TargetID: req.CommentId, UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return &postpb.LikeCountResponse{LikesCount: likes}, nil
}

func (s *PostServer) UnlikeComment(ctx context.Context, req *postpb.UnlikeCommentRequest) (*postpb.LikeCountResponse, error) {
	likes, err := s.services.Likes.UnlikeComment(ctx, models.LikeInput{TargetID: req.CommentId, UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	return &postpb.LikeCountResponse{LikesCount: likes}, nil
}

func (s *PostServer) ListLabels(ctx context.Context, _ *postpb.Empty) (*postpb.ListLabelsResponse, error) {
	labels, err := s.services.Labels.ListLabels(ctx)
	if err != nil {
		return nil, err
	}
	return &postpb.ListLabelsResponse{Labels: toProtoLabels(labels)}, nil
}

func (s *PostServer) AddLabelToPost(ctx context.Context, req *postpb.AddLabelToPostRequest) (*postpb.Empty, error) {
	err := s.services.Labels.AddLabelToPost(ctx, models.LabelAssignmentInput{PostID: req.PostId, LabelID: req.LabelId})
	if err != nil {
		return nil, err
	}
	return &postpb.Empty{}, nil
}

func (s *PostServer) RemoveLabelFromPost(ctx context.Context, req *postpb.RemoveLabelFromPostRequest) (*postpb.Empty, error) {
	err := s.services.Labels.RemoveLabelFromPost(ctx, models.LabelAssignmentInput{PostID: req.PostId, LabelID: req.LabelId})
	if err != nil {
		return nil, err
	}
	return &postpb.Empty{}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.Post, error) {
	var titlePtr *string
	if req.Title != "" {
		title := req.Title
		titlePtr = &title
	}
	var descriptionPtr *string
	if req.Description != "" {
		description := req.Description
		descriptionPtr = &description
	}
	post, err := s.services.Posts.UpdatePost(ctx, models.UpdatePostInput{ID: req.Id, Title: titlePtr, Description: descriptionPtr})
	if err != nil {
		return nil, err
	}
	return toProtoPost(*post), nil
}

func (s *PostServer) DeletePost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.Empty, error) {
	if err := s.services.Posts.DeletePost(ctx, req.Id); err != nil {
		return nil, err
	}
	return &postpb.Empty{}, nil
}

func toProtoPost(post models.Post) *postpb.Post {
	return &postpb.Post{
		Id:          post.ID,
		UserId:      post.UserID,
		Title:       post.Title,
		Description: post.Description,
		LikesCount:  post.LikesCount,
		CreatedAt:   post.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   post.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Labels:      toProtoLabels(post.Labels),
	}
}

func toProtoComment(comment models.Comment) *postpb.Comment {
	return &postpb.Comment{
		Id:         comment.ID,
		PostId:     comment.PostID,
		UserId:     comment.UserID,
		Content:    comment.Content,
		LikesCount: comment.LikesCount,
		CreatedAt:  comment.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:  comment.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func toProtoComments(comments []models.Comment) []*postpb.Comment {
	proto := make([]*postpb.Comment, 0, len(comments))
	for _, comment := range comments {
		proto = append(proto, toProtoComment(comment))
	}
	return proto
}

func toProtoLabels(labels []models.Label) []*postpb.Label {
	proto := make([]*postpb.Label, 0, len(labels))
	for _, label := range labels {
		proto = append(proto, &postpb.Label{Id: label.ID, Name: label.Name})
	}
	return proto
}
