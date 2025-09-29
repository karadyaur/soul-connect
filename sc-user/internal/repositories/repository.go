package repositories

import db "soul-connect/sc-user/internal/db/sqlc"

type Repository struct {
	User         IUserRepository
	Subscription ISubscriptionRepository
}

func NewRepository(queries db.Querier) *Repository {
	return &Repository{
		User:         NewUserRepository(queries),
		Subscription: NewSubscriptionRepository(queries),
	}
}
