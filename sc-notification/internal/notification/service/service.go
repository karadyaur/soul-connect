package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	db "soul-connect/sc-notification/internal/db/sqlc"
	"soul-connect/sc-notification/internal/notification/repository"
)

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationService interface {
	CreateNotification(ctx context.Context, userID uuid.UUID, content string) (Notification, error)
	ListNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error)
}

type service struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &service{repo: repo}
}

func (s *service) CreateNotification(ctx context.Context, userID uuid.UUID, content string) (Notification, error) {
	record, err := s.repo.CreateNotification(ctx, userID, content)
	if err != nil {
		return Notification{}, err
	}

	return mapNotification(record)
}

func (s *service) ListNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error) {
	records, err := s.repo.ListNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	notifications := make([]Notification, 0, len(records))
	for _, record := range records {
		notification, err := mapNotification(record)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func mapNotification(record db.Notification) (Notification, error) {
	if !record.ID.Valid {
		return Notification{}, errors.New("notification id is invalid")
	}
	if !record.UserID.Valid {
		return Notification{}, errors.New("notification user id is invalid")
	}
	if !record.CreatedAt.Valid {
		return Notification{}, errors.New("notification creation time is invalid")
	}

	id := uuid.UUID(record.ID.Bytes)
	userID := uuid.UUID(record.UserID.Bytes)

	return Notification{
		ID:        id.String(),
		UserID:    userID.String(),
		Content:   record.Content,
		CreatedAt: record.CreatedAt.Time,
	}, nil
}
