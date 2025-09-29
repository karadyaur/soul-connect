package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	db "soul-connect/sc-user/internal/db/sqlc"
	"soul-connect/sc-user/internal/repositories"
)

type Service struct {
	UserService *UserService
}

func NewService(pool *pgxpool.Pool) *Service {
	queries := db.New(pool)
	repo := repositories.NewRepository(queries)
	return &Service{
		UserService: NewUserService(repo.User, repo.Subscription),
	}
}
