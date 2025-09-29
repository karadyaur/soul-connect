package repositories

import (
	"context"

	db "soul-connect/sc-user/internal/db/sqlc"
	"soul-connect/sc-user/internal/models"
)

type IUserRepository interface {
	Create(ctx context.Context, input models.CreateUserProfileInput) (*models.UserProfile, error)
	GetByID(ctx context.Context, id string) (*models.UserProfile, error)
	Update(ctx context.Context, params models.UpdateUserProfileParams) (*models.UserProfile, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository struct {
	queries db.Querier
}

func NewUserRepository(queries db.Querier) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, input models.CreateUserProfileInput) (*models.UserProfile, error) {
	authID, err := stringToUUID(input.AuthID)
	if err != nil {
		return nil, err
	}

	created, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		AuthID:    authID,
		FullName:  input.FullName,
		Bio:       stringPtrToText(input.Bio),
		PhotoLink: stringPtrToText(input.PhotoLink),
	})
	if err != nil {
		return nil, err
	}

	return userToModel(created), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.UserProfile, error) {
	userID, err := stringToUUID(id)
	if err != nil {
		return nil, err
	}

	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userToModel(user), nil
}

func (r *UserRepository) Update(ctx context.Context, params models.UpdateUserProfileParams) (*models.UserProfile, error) {
	userID, err := stringToUUID(params.ID)
	if err != nil {
		return nil, err
	}

	if err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        userID,
		FullName:  params.FullName,
		Bio:       stringPtrToText(params.Bio),
		PhotoLink: stringPtrToText(params.PhotoLink),
	}); err != nil {
		return nil, err
	}

	updated, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userToModel(updated), nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	userID, err := stringToUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteUser(ctx, userID)
}

var _ IUserRepository = (*UserRepository)(nil)
