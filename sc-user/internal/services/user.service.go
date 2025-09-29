package services

import (
	"context"
	"errors"

	"soul-connect/sc-user/internal/models"
	"soul-connect/sc-user/internal/repositories"
)

type UserService struct {
	userRepo         repositories.IUserRepository
	subscriptionRepo repositories.ISubscriptionRepository
}

func NewUserService(userRepo repositories.IUserRepository, subscriptionRepo repositories.ISubscriptionRepository) *UserService {
	return &UserService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *UserService) CreateProfile(ctx context.Context, input models.CreateUserProfileInput) (*models.UserProfile, error) {
	if input.FullName == "" {
		return nil, errors.New("full name is required")
	}
	if input.AuthID == "" {
		return nil, errors.New("auth id is required")
	}
	return s.userRepo.Create(ctx, input)
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, input models.UpdateUserProfileInput) (*models.UserProfile, error) {
	if input.ID == "" {
		return nil, errors.New("user id is required")
	}

	existing, err := s.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	fullName := existing.FullName
	if input.FullName != nil {
		fullName = *input.FullName
	}

	params := models.UpdateUserProfileParams{
		ID:        input.ID,
		FullName:  fullName,
		Bio:       input.Bio,
		PhotoLink: input.PhotoLink,
	}

	return s.userRepo.Update(ctx, params)
}

func (s *UserService) DeleteProfile(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user id is required")
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) Subscribe(ctx context.Context, input models.ModifySubscriptionInput) error {
	if input.SubscriberID == "" || input.AuthorID == "" {
		return errors.New("subscriber and author ids are required")
	}
	if input.SubscriberID == input.AuthorID {
		return errors.New("subscriber and author cannot be the same user")
	}
	return s.subscriptionRepo.Subscribe(ctx, input.SubscriberID, input.AuthorID)
}

func (s *UserService) Unsubscribe(ctx context.Context, input models.ModifySubscriptionInput) error {
	if input.SubscriberID == "" || input.AuthorID == "" {
		return errors.New("subscriber and author ids are required")
	}
	return s.subscriptionRepo.Unsubscribe(ctx, input.SubscriberID, input.AuthorID)
}

func (s *UserService) ListSubscriptions(ctx context.Context, subscriberID string) (*models.SubscriptionList, error) {
	if subscriberID == "" {
		return nil, errors.New("subscriber id is required")
	}

	authorIDs, err := s.subscriptionRepo.ListAuthorIDs(ctx, subscriberID)
	if err != nil {
		return nil, err
	}

	return &models.SubscriptionList{
		SubscriberID: subscriberID,
		AuthorIDs:    authorIDs,
	}, nil
}
