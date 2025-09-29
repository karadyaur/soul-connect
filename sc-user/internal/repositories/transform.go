package repositories

import (
	"time"

	"github.com/google/uuid"
	db "soul-connect/sc-user/internal/db/sqlc"
	"soul-connect/sc-user/internal/models"
)

func userToModel(user db.User) *models.UserProfile {
	var id string
	if user.ID.Valid {
		parsed := uuid.UUID(user.ID.Bytes)
		id = parsed.String()
	}
	var authID string
	if user.AuthID.Valid {
		parsed := uuid.UUID(user.AuthID.Bytes)
		authID = parsed.String()
	}

	var createdAt time.Time
	if user.CreatedAt.Valid {
		createdAt = user.CreatedAt.Time
	}
	var updatedAt time.Time
	if user.UpdatedAt.Valid {
		updatedAt = user.UpdatedAt.Time
	}

	return &models.UserProfile{
		ID:        id,
		AuthID:    authID,
		FullName:  user.FullName,
		Bio:       textToStringPtr(user.Bio),
		PhotoLink: textToStringPtr(user.PhotoLink),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
