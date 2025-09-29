package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "soul-connect/sc-notification/internal/db/sqlc"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, userID uuid.UUID, content string) (db.Notification, error)
	ListNotifications(ctx context.Context, userID uuid.UUID) ([]db.Notification, error)
}

type repository struct {
	queries db.Querier
}

func NewNotificationRepository(queries db.Querier) NotificationRepository {
	return &repository{queries: queries}
}

func (r *repository) CreateNotification(ctx context.Context, userID uuid.UUID, content string) (db.Notification, error) {
	pgUserID := pgtype.UUID{Bytes: userID, Valid: true}
	params := db.CreateNotificationParams{
		UserID:  pgUserID,
		Content: content,
	}

	notification, err := r.queries.CreateNotification(ctx, params)
	if err != nil {
		return db.Notification{}, fmt.Errorf("create notification: %w", err)
	}

	return notification, nil
}

func (r *repository) ListNotifications(ctx context.Context, userID uuid.UUID) ([]db.Notification, error) {
	pgUserID := pgtype.UUID{Bytes: userID, Valid: true}

	notifications, err := r.queries.GetNotificationsByUser(ctx, pgUserID)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}

	return notifications, nil
}
