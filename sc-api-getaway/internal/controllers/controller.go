package controllers

import "soul-connect/sc-api-getaway/internal/generated"

type Controller struct {
	AuthController *AuthController
	UserController *UserController
}

func NewController(authClient generated.AuthServiceClient, userClient generated.UserServiceClient) *Controller {
	return &Controller{
		AuthController: NewAuthController(authClient),
		UserController: NewUserController(userClient),
	}
}
