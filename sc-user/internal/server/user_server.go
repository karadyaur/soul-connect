package server

import (
	"context"
	"time"

	"soul-connect/sc-user/internal/generated"
	"soul-connect/sc-user/internal/models"
	"soul-connect/sc-user/internal/services"
)

type UserServer struct {
	generated.UnimplementedUserServiceServer
	userService *services.UserService
}

func NewUserServer(service *services.Service) *UserServer {
	return &UserServer{
		userService: service.UserService,
	}
}

func (s *UserServer) CreateProfile(ctx context.Context, request *generated.CreateProfileRequest) (*generated.UserProfile, error) {
	input := models.CreateUserProfileInput{
		AuthID:   request.AuthId,
		FullName: request.FullName,
	}
	if request.Bio != nil {
		bio := request.GetBio()
		input.Bio = &bio
	}
	if request.PhotoLink != nil {
		photo := request.GetPhotoLink()
		input.PhotoLink = &photo
	}

	profile, err := s.userService.CreateProfile(ctx, input)
	if err != nil {
		return nil, err
	}

	return toProtoProfile(profile), nil
}

func (s *UserServer) GetProfile(ctx context.Context, request *generated.GetProfileRequest) (*generated.UserProfile, error) {
	profile, err := s.userService.GetProfile(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return toProtoProfile(profile), nil
}

func (s *UserServer) UpdateProfile(ctx context.Context, request *generated.UpdateProfileRequest) (*generated.UserProfile, error) {
	input := models.UpdateUserProfileInput{ID: request.Id}
	if request.FullName != nil {
		value := request.GetFullName()
		input.FullName = &value
	}
	if request.Bio != nil {
		value := request.GetBio()
		input.Bio = &value
	}
	if request.PhotoLink != nil {
		value := request.GetPhotoLink()
		input.PhotoLink = &value
	}

	profile, err := s.userService.UpdateProfile(ctx, input)
	if err != nil {
		return nil, err
	}
	return toProtoProfile(profile), nil
}

func (s *UserServer) DeleteProfile(ctx context.Context, request *generated.DeleteProfileRequest) (*generated.Empty, error) {
	if err := s.userService.DeleteProfile(ctx, request.Id); err != nil {
		return nil, err
	}
	return &generated.Empty{}, nil
}

func (s *UserServer) Subscribe(ctx context.Context, request *generated.ModifySubscriptionRequest) (*generated.Empty, error) {
	input := models.ModifySubscriptionInput{
		SubscriberID: request.SubscriberId,
		AuthorID:     request.AuthorId,
	}
	if err := s.userService.Subscribe(ctx, input); err != nil {
		return nil, err
	}
	return &generated.Empty{}, nil
}

func (s *UserServer) Unsubscribe(ctx context.Context, request *generated.ModifySubscriptionRequest) (*generated.Empty, error) {
	input := models.ModifySubscriptionInput{
		SubscriberID: request.SubscriberId,
		AuthorID:     request.AuthorId,
	}
	if err := s.userService.Unsubscribe(ctx, input); err != nil {
		return nil, err
	}
	return &generated.Empty{}, nil
}

func (s *UserServer) ListSubscriptions(ctx context.Context, request *generated.ListSubscriptionsRequest) (*generated.ListSubscriptionsResponse, error) {
	result, err := s.userService.ListSubscriptions(ctx, request.SubscriberId)
	if err != nil {
		return nil, err
	}
	return &generated.ListSubscriptionsResponse{
		SubscriberId: result.SubscriberID,
		AuthorIds:    result.AuthorIDs,
	}, nil
}

func toProtoProfile(profile *models.UserProfile) *generated.UserProfile {
	if profile == nil {
		return nil
	}
	response := &generated.UserProfile{
		Id:        profile.ID,
		AuthId:    profile.AuthID,
		FullName:  profile.FullName,
		CreatedAt: profile.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339Nano),
	}
	if profile.Bio != nil {
		bio := *profile.Bio
		response.Bio = &bio
	}
	if profile.PhotoLink != nil {
		link := *profile.PhotoLink
		response.PhotoLink = &link
	}
	return response
}
