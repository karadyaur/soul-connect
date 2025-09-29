package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"soul-connect/sc-api-getaway/internal/generated"
)

type UserController struct {
	client generated.UserServiceClient
}

func NewUserController(userClient generated.UserServiceClient) *UserController {
	return &UserController{client: userClient}
}

type createProfilePayload struct {
	AuthID    string  `json:"auth_id"`
	FullName  string  `json:"full_name"`
	Bio       *string `json:"bio"`
	PhotoLink *string `json:"photo_link"`
}

type updateProfilePayload struct {
	FullName  *string `json:"full_name"`
	Bio       *string `json:"bio"`
	PhotoLink *string `json:"photo_link"`
}

type subscriptionPayload struct {
	AuthorID string `json:"author_id"`
}

func (c *UserController) CreateProfile(gc *gin.Context) {
	var payload createProfilePayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if payload.AuthID == "" || payload.FullName == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "auth_id and full_name are required"})
		return
	}

	req := &generated.CreateProfileRequest{
		AuthId:   payload.AuthID,
		FullName: payload.FullName,
	}
	if payload.Bio != nil {
		req.Bio = payload.Bio
	}
	if payload.PhotoLink != nil {
		req.PhotoLink = payload.PhotoLink
	}

	ctx := context.Background()
	profile, err := c.client.CreateProfile(ctx, req)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusCreated, toUserProfileResponse(profile))
}

func (c *UserController) GetProfile(gc *gin.Context) {
	id := gc.Param("id")
	if id == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ctx := context.Background()
	profile, err := c.client.GetProfile(ctx, &generated.GetProfileRequest{Id: id})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, toUserProfileResponse(profile))
}

func (c *UserController) UpdateProfile(gc *gin.Context) {
	id := gc.Param("id")
	if id == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var payload updateProfilePayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	req := &generated.UpdateProfileRequest{Id: id}
	if payload.FullName != nil {
		req.FullName = payload.FullName
	}
	if payload.Bio != nil {
		req.Bio = payload.Bio
	}
	if payload.PhotoLink != nil {
		req.PhotoLink = payload.PhotoLink
	}

	ctx := context.Background()
	profile, err := c.client.UpdateProfile(ctx, req)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, toUserProfileResponse(profile))
}

func (c *UserController) DeleteProfile(gc *gin.Context) {
	id := gc.Param("id")
	if id == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ctx := context.Background()
	if _, err := c.client.DeleteProfile(ctx, &generated.DeleteProfileRequest{Id: id}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.Status(http.StatusNoContent)
}

func (c *UserController) Subscribe(gc *gin.Context) {
	subscriberID := gc.Param("id")
	if subscriberID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "subscriber id is required"})
		return
	}

	var payload subscriptionPayload
	if err := gc.ShouldBindJSON(&payload); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if payload.AuthorID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "author_id is required"})
		return
	}

	ctx := context.Background()
	if _, err := c.client.Subscribe(ctx, &generated.ModifySubscriptionRequest{
		SubscriberId: subscriberID,
		AuthorId:     payload.AuthorID,
	}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.Status(http.StatusCreated)
}

func (c *UserController) Unsubscribe(gc *gin.Context) {
	subscriberID := gc.Param("id")
	authorID := gc.Param("author_id")
	if subscriberID == "" || authorID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "subscriber and author ids are required"})
		return
	}

	ctx := context.Background()
	if _, err := c.client.Unsubscribe(ctx, &generated.ModifySubscriptionRequest{
		SubscriberId: subscriberID,
		AuthorId:     authorID,
	}); err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.Status(http.StatusNoContent)
}

func (c *UserController) ListSubscriptions(gc *gin.Context) {
	subscriberID := gc.Param("id")
	if subscriberID == "" {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "subscriber id is required"})
		return
	}

	ctx := context.Background()
	resp, err := c.client.ListSubscriptions(ctx, &generated.ListSubscriptionsRequest{SubscriberId: subscriberID})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"subscriber_id": resp.SubscriberId,
		"author_ids":    resp.AuthorIds,
	})
}

func toUserProfileResponse(profile *generated.UserProfile) gin.H {
	if profile == nil {
		return gin.H{}
	}
	response := gin.H{
		"id":         profile.Id,
		"auth_id":    profile.AuthId,
		"full_name":  profile.FullName,
		"created_at": profile.CreatedAt,
		"updated_at": profile.UpdatedAt,
	}
	if profile.Bio != nil {
		response["bio"] = profile.GetBio()
	}
	if profile.PhotoLink != nil {
		response["photo_link"] = profile.GetPhotoLink()
	}
	return response
}
