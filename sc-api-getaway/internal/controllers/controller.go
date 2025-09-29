package controllers

import (
	"soul-connect/sc-api-getaway/internal/generated"
	postpb "soul-connect/sc-post/pkg/postpb"
)

type Controller struct {
	AuthController *AuthController
	PostController *PostController
}

func NewController(authClient generated.AuthServiceClient, postClient postpb.PostServiceClient) *Controller {
	return &Controller{
		AuthController: NewAuthController(authClient),
		PostController: NewPostController(postClient),
	}
}
