package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	postpb "soul-connect/sc-post/pkg/postpb"
)

type PostController struct {
	client postpb.PostServiceClient
}

func NewPostController(client postpb.PostServiceClient) *PostController {
	return &PostController{client: client}
}

func (c *PostController) CreatePost(gc *gin.Context) {
	var req struct {
		UserID      string   `json:"user_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		LabelIDs    []string `json:"label_ids"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.UserID == "" || req.Title == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "user_id and title are required"})
		return
	}

	ctx := context.Background()
	resp, err := c.client.CreatePost(ctx, &postpb.CreatePostRequest{
		UserId:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		LabelIds:    req.LabelIDs,
	})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusCreated, postToResponse(resp.Post))
}

func (c *PostController) GetPost(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	ctx := context.Background()
	resp, err := c.client.GetPost(ctx, &postpb.GetPostRequest{Id: postID})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{
		"post":     postToResponse(resp.Post),
		"comments": commentsToResponse(resp.Comments),
	})
}

func (c *PostController) ListPosts(gc *gin.Context) {
	var labelIDs []string
	labelsParam := gc.Query("labels")
	if labelsParam != "" {
		for _, part := range strings.Split(labelsParam, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				labelIDs = append(labelIDs, trimmed)
			}
		}
	}
	ctx := context.Background()
	resp, err := c.client.ListPosts(ctx, &postpb.ListPostsRequest{LabelIds: labelIDs})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	posts := make([]gin.H, 0, len(resp.Posts))
	for _, summary := range resp.Posts {
		posts = append(posts, gin.H{
			"post":           postToResponse(summary.Post),
			"total_comments": summary.TotalComments,
			"total_likes":    summary.TotalLikes,
		})
	}
	gc.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (c *PostController) AddComment(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	var req struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.UserID == "" || req.Content == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "user_id and content are required"})
		return
	}
	ctx := context.Background()
	resp, err := c.client.AddComment(ctx, &postpb.AddCommentRequest{PostId: postID, UserId: req.UserID, Content: req.Content})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusCreated, commentToResponse(resp.Comment))
}

func (c *PostController) ListComments(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	ctx := context.Background()
	resp, err := c.client.ListComments(ctx, &postpb.ListCommentsRequest{PostId: postID})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{"comments": commentsToResponse(resp.Comments)})
}

func (c *PostController) LikePost(gc *gin.Context) {
	c.handleLike(gc, true)
}

func (c *PostController) UnlikePost(gc *gin.Context) {
	c.handleLike(gc, false)
}

func (c *PostController) LikeComment(gc *gin.Context) {
	c.handleCommentLike(gc, true)
}

func (c *PostController) UnlikeComment(gc *gin.Context) {
	c.handleCommentLike(gc, false)
}

func (c *PostController) ListLabels(gc *gin.Context) {
	ctx := context.Background()
	resp, err := c.client.ListLabels(ctx, &postpb.Empty{})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{"labels": labelsToResponse(resp.Labels)})
}

func (c *PostController) AddLabelToPost(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	var req struct {
		LabelID string `json:"label_id"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil || req.LabelID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "label_id is required"})
		return
	}
	ctx := context.Background()
	if _, err := c.client.AddLabelToPost(ctx, &postpb.AddLabelToPostRequest{PostId: postID, LabelId: req.LabelID}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.Status(http.StatusNoContent)
}

func (c *PostController) RemoveLabelFromPost(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	labelID := gc.Param("label_id")
	if labelID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "label_id is required"})
		return
	}
	ctx := context.Background()
	if _, err := c.client.RemoveLabelFromPost(ctx, &postpb.RemoveLabelFromPostRequest{PostId: postID, LabelId: labelID}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.Status(http.StatusNoContent)
}

func (c *PostController) UpdatePost(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	ctx := context.Background()
	resp, err := c.client.UpdatePost(ctx, &postpb.UpdatePostRequest{Id: postID, Title: req.Title, Description: req.Description})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, postToResponse(resp))
}

func (c *PostController) DeletePost(gc *gin.Context) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	ctx := context.Background()
	if _, err := c.client.DeletePost(ctx, &postpb.GetPostRequest{Id: postID}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.Status(http.StatusNoContent)
}

func (c *PostController) handleLike(gc *gin.Context, like bool) {
	postID := gc.Param("post_id")
	if postID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	ctx := context.Background()
	var (
		resp *postpb.LikeCountResponse
		err  error
	)
	if like {
		resp, err = c.client.LikePost(ctx, &postpb.LikePostRequest{PostId: postID, UserId: req.UserID})
	} else {
		resp, err = c.client.UnlikePost(ctx, &postpb.UnlikePostRequest{PostId: postID, UserId: req.UserID})
	}
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{"likes_count": resp.LikesCount})
}

func (c *PostController) handleCommentLike(gc *gin.Context, like bool) {
	commentID := gc.Param("comment_id")
	if commentID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "comment_id is required"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := gc.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	ctx := context.Background()
	var (
		resp *postpb.LikeCountResponse
		err  error
	)
	if like {
		resp, err = c.client.LikeComment(ctx, &postpb.LikeCommentRequest{CommentId: commentID, UserId: req.UserID})
	} else {
		resp, err = c.client.UnlikeComment(ctx, &postpb.UnlikeCommentRequest{CommentId: commentID, UserId: req.UserID})
	}
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.JSON(http.StatusOK, gin.H{"likes_count": resp.LikesCount})
}

func postToResponse(post *postpb.Post) gin.H {
	if post == nil {
		return gin.H{}
	}
	return gin.H{
		"id":          post.Id,
		"user_id":     post.UserId,
		"title":       post.Title,
		"description": post.Description,
		"likes_count": post.LikesCount,
		"created_at":  post.CreatedAt,
		"updated_at":  post.UpdatedAt,
		"labels":      labelsToResponse(post.Labels),
	}
}

func commentToResponse(comment *postpb.Comment) gin.H {
	if comment == nil {
		return gin.H{}
	}
	return gin.H{
		"id":          comment.Id,
		"post_id":     comment.PostId,
		"user_id":     comment.UserId,
		"content":     comment.Content,
		"likes_count": comment.LikesCount,
		"created_at":  comment.CreatedAt,
		"updated_at":  comment.UpdatedAt,
	}
}

func commentsToResponse(comments []*postpb.Comment) []gin.H {
	result := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		result = append(result, commentToResponse(comment))
	}
	return result
}

func labelsToResponse(labels []*postpb.Label) []gin.H {
	result := make([]gin.H, 0, len(labels))
	for _, label := range labels {
		result = append(result, gin.H{"id": label.Id, "name": label.Name})
	}
	return result
}
